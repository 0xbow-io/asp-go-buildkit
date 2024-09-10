package main

import (
	"context"
	"fmt"
	"math/big"

	sql "github.com/iden3/go-merkletree-sql/db/pgx/v5" //note this import path, official example is wrong
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// pool is used instead of conn as its concurrency safe
func setupPostgresConnection() (*pgxpool.Pool, error) {
	connString := "postgres://myuser:mypassword@localhost:5432/myproject"
	return pgxpool.New(context.Background(), connString)
}

func createSQLStorage(pool *pgxpool.Pool) merkletree.Storage {
	return sql.NewSqlStorage(pool, 1) // 1 is the tree ID
}

// creates the necessary tables in postgres (required)
func initializeDatabase(pool *pgxpool.Pool) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS mt_roots (
            mt_id BIGINT NOT NULL,
            key BYTEA NOT NULL,
            created_at BIGINT,
            deleted_at BIGINT,
            PRIMARY KEY (mt_id)
        );

        CREATE TABLE IF NOT EXISTS mt_nodes (
            mt_id BIGINT NOT NULL,
            key BYTEA NOT NULL,
            type SMALLINT NOT NULL,
            child_l BYTEA,
            child_r BYTEA,
            entry BYTEA,
            created_at BIGINT,
            deleted_at BIGINT,
            PRIMARY KEY (mt_id, key)
        );
    `)
	return err
}

func main() {
	ctx := context.Background()

	// Set up PostgreSQL connection
	dbPool, err := setupPostgresConnection()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to PostgreSQL: %v", err))
	}
	defer dbPool.Close()

	// Initialize the database schema
	err = initializeDatabase(dbPool)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}

	// Create SQL storage
	storage := createSQLStorage(dbPool)

	// Create a new MerkleTree with a depth of 10
	mt, err := merkletree.NewMerkleTree(ctx, storage, 10)
	if err != nil {
		panic(fmt.Sprintf("Failed to create MerkleTree: %v", err))
	}

	// Add some elements to the tree
	err = mt.Add(ctx, big.NewInt(1), big.NewInt(10))
	if err != nil {
		panic(fmt.Sprintf("Failed to add element: %v", err))
	}
	err = mt.Add(ctx, big.NewInt(2), big.NewInt(20))
	if err != nil {
		panic(fmt.Sprintf("Failed to add element: %v", err))
	}

	// Generate a proof for an element
	proof, value, err := mt.GenerateProof(ctx, big.NewInt(1), nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate proof: %v", err))
	}

	// Verify the proof
	valid := merkletree.VerifyProof(mt.Root(), proof, big.NewInt(1), value)
	fmt.Printf("Proof verification result: %v\n", valid)

	// Print the root of the tree
	fmt.Printf("Merkle Tree Root: %s\n", mt.Root().String())
}
