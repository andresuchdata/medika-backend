package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"medika-backend/internal/infrastructure/config"
	"medika-backend/internal/infrastructure/database"
	"medika-backend/internal/infrastructure/seeder"
)

func main() {
	var (
		env = flag.String("env", "development", "Environment (development, production)")
		all = flag.Bool("all", false, "Run all seeders")
		org = flag.Bool("organizations", false, "Seed organizations")
		usr = flag.Bool("users", false, "Seed users")
		rom = flag.Bool("rooms", false, "Seed rooms")
	)
	flag.Parse()

	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create seeder manager
	manager := seeder.NewManager()

	// Register seeders based on flags
	if *all || *org {
		manager.Register(seeder.NewOrganizationSeeder())
	}
	if *all || *usr {
		manager.Register(seeder.NewUserSeeder())
	}
	if *all || *rom {
		manager.Register(seeder.NewRoomSeeder())
	}

	// Check if any seeders were registered
	if len(manager.GetSeeders()) == 0 {
		fmt.Println("Usage: seeder [options]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  seeder -all                    # Run all seeders")
		fmt.Println("  seeder -organizations -users   # Seed organizations and users")
		fmt.Println("  seeder -users                  # Seed only users")
		os.Exit(1)
	}

	// Confirmation for production
	if *env == "production" {
		fmt.Print("⚠️  You are about to seed data in PRODUCTION environment. Continue? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" && response != "yes" {
			fmt.Println("Seeding cancelled.")
			os.Exit(0)
		}
	}

	// Run seeders
	fmt.Printf("🌱 Starting seeding process in %s environment...\n\n", *env)
	
	for _, s := range manager.GetSeeders() {
		fmt.Printf("Running %s...\n", s.Name())
		if err := s.Seed(ctx, db); err != nil {
			log.Fatalf("❌ Seeder %s failed: %v", s.Name(), err)
		}
	}

	fmt.Println("\n🎉 All seeders completed successfully!")
}
