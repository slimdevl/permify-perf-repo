package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/slimdevl/repro/pkg/permify"
)

const (
	TenantId          = "test"
	DefaultIterations = 100
	DefaultCount      = 100
	MaxSleep          = 100
)

func main() {
	var maxIterations int
	var relationCount int
	var rateLimit int

	// Define command-line flags
	flag.IntVar(&maxIterations, "iterations", DefaultIterations, "Number of iterations")
	flag.IntVar(&relationCount, "count", DefaultCount, "Number of iterations")
	flag.IntVar(&rateLimit, "rate-limit", permify.DefaultRateLimit, "Rate limit")

	// Parse the command-line flags
	flag.Parse()

	ctx := context.Background()
	cfg := permify.NewDefaultConfig()
	cfg.Tenant = TenantId
	cfg.RateLimit = rateLimit
	client := permify.NewClient(cfg)

	// Create a new tenant
	_, err := client.(permify.SchemaManagerClient).CreateTenant(ctx, &permify.CreateTenantRequest{})
	if err != nil {
		log.Fatalf("Error creating tenant: %v\n", err)
	}
	_, err = client.(permify.SchemaManagerClient).SaveModelSchema(ctx, &permify.SaveSchemaRequest{
		Schema: testSchema,
	})
	if err != nil {
		log.Fatalf("Error saving schema: %v\n", err)
	}

	fmt.Printf("Beginning test to ensure we can create and destroy the same relationships\n")
	fmt.Printf("iterations: %d, count: %d\n", maxIterations, relationCount)

	cleanupCh := make(chan int, relationCount)

	for i := 1; i <= maxIterations; i++ {
		wg := sync.WaitGroup{}
		fmt.Printf("iterations %d starting", i)
		relationshipSets := relationshipGenerator(0, relationCount)

		// Start goroutines for delete operations
		for i := range relationshipSets {
			wg.Add(1)
			go func(wg *sync.WaitGroup, entry int) {
				defer wg.Done()
				index := <-cleanupCh
				set := relationshipSets[index]
				deleteRelationships(client, ctx, set)
				fmt.Printf("x")
			}(&wg, i)
		}

		// Start goroutines for adding relationships
		for i := 0; i < len(relationshipSets); i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, index int) {
				defer wg.Done()
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(MaxSleep+1)))
				set := relationshipSets[index]
				addRelationships(client, ctx, set)
				cleanupCh <- index
				fmt.Printf(".")
			}(&wg, i)
		}

		wg.Wait()
		fmt.Println("done")
	}

	close(cleanupCh)
	fmt.Println("testing complete")
}

func addRelationships(client permify.RelationshipClient, ctx context.Context, relationships []*permify.Relationship) {
	if len(relationships) == 0 {
		log.Fatalf("No relationships to add in addRelationships! eek\n")
	}
	if _, err := client.AddRelationship(ctx, &permify.AddRelationshipRequest{
		Relationships: relationships,
	}); err != nil {
		log.Fatalf("Error adding relationship: %v\n", err)
	}
}

func deleteRelationships(client permify.RelationshipClient, ctx context.Context, relationships []*permify.Relationship) {
	// now delete the relationship in reverse order
	for i := len(relationships) - 1; i >= 0; i-- {
		relation := relationships[i]
		err := client.DeleteRelationship(ctx, &permify.DeleteRelationshipRequest{
			Filter: permify.RelationshipFilter{
				Entity: permify.EntityIDSet{
					Type: relation.Entity.Type,
					Ids:  []string{relation.Entity.Id},
				},
				Relation: relation.Relation,
				Subject: permify.SubjectIDSet{
					Type: relation.Subject.Type,
					Ids:  []string{relation.Subject.Id},
				},
			},
		})
		if err != nil {
			log.Fatalf("Error deleting relationship %s:%s -> %s -> %s:%s: %s\n",
				relation.Entity.Type, relation.Entity.Id,
				relation.Relation,
				relation.Subject.Type, relation.Subject.Id,
				err.Error())
		}
	}
}

const (
	User         = "user"
	Organization = "organization"
	Team         = "team"
	Project      = "project"
)

const (
	RoleOrg    = "org"
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"
	RoleTeam   = "team"
)

// Tests schema from Permify examples
const testSchema = `
entity user {}

entity organization {

	// organizational roles
	relation admin @user
	relation member @user

}

entity team {
	// reference for organization that team belong
	relation org @organization

	// represents owner or creator of the team
	relation owner @user

	// represents direct member of the team
	relation member @user

	// organization admins or owners can edit, delete the team details
	permission edit = org.admin or owner
	permission delete = org.admin or owner

	// to invite someone you need to be admin and either owner or member of this team
	permission invite = org.admin and (owner or member)

	// only owners can remove users
	permission remove_user =  owner
}

entity project {

	// references for team and organization that project belongs
	relation team @team
	relation org @organization

	permission view = org.admin or team.member
	permission edit = org.admin or team.member
	permission delete = team.member
}
`

// golang generator which will create n number of relationship set
// starting at a specific value
func relationshipGenerator(start, n int) [][]*permify.Relationship {
	relationshipSets := make([][]*permify.Relationship, n)
	for i := 0; i < n; i++ {
		relationshipSets[i] = makeRelationships(start + i)
	}
	return relationshipSets
}

func ID(t string, i int) string {
	return fmt.Sprintf("%s.%d", t, i)
}

func makeRelationships(i int) []*permify.Relationship {
	return []*permify.Relationship{
		// add user to org as admin
		{
			Entity:   &permify.Entity{Type: Organization, Id: ID(Organization, i)},
			Relation: RoleAdmin,
			Subject:  &permify.Subject{Type: User, Id: ID(User, i)},
		},
		// add user to org as member
		{
			Entity:   &permify.Entity{Type: Organization, Id: ID(Organization, i)},
			Relation: RoleMember,
			Subject:  &permify.Subject{Type: User, Id: ID(User, i)},
		},
		// add organization to team as org
		{
			Entity:   &permify.Entity{Type: Team, Id: ID(Team, i)},
			Relation: RoleOrg,
			Subject:  &permify.Subject{Type: Organization, Id: ID(Organization, i)},
		},
		// add user to team owner
		{
			Entity:   &permify.Entity{Type: Team, Id: ID(Team, i)},
			Relation: RoleOwner,
			Subject:  &permify.Subject{Type: User, Id: ID(User, i)},
		},
		// add user to team member
		{
			Entity:   &permify.Entity{Type: Team, Id: ID(Team, i)},
			Relation: RoleMember,
			Subject:  &permify.Subject{Type: User, Id: ID(User, i)},
		},
		// add organization to project as team
		{
			Entity:   &permify.Entity{Type: Project, Id: ID(Project, i)},
			Relation: RoleOrg,
			Subject:  &permify.Subject{Type: Organization, Id: ID(Organization, i)},
		},
		// add team to project as team
		{
			Entity:   &permify.Entity{Type: Project, Id: ID(Project, i)},
			Relation: RoleTeam,
			Subject:  &permify.Subject{Type: Team, Id: ID(Team, i)},
		},
	}
}
