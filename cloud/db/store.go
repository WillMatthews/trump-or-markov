package db

type StoreFactory interface {
	NewStore() Store
}

type Store interface {
}

// Plan:
// 1 -> setup sqlite as the store to use
// 2 -> create a factory for the store
// 3 -> store tweets / other raw data in the store (?)
// 4 -> store the generated markov chains & dictionaries in the store
// 5 -> store user data in the store (? do I want to do this, maybe, for the sake of having a high score list)
// 6 -> store 'favourited' fake tweets in the store
// 7 -> real tweet search?
//
// To do all of the above I need to write a schema.
//
// Tech stack:
// - (db) https://github.com/electric-sql/pglite (if that fails, just use sqlite) https://github.com/fergusstrange/embedded-postgres
// - (pg driver) pgx with scany (if using sqlite, then use sqlx)
// - (migrations) goose?
