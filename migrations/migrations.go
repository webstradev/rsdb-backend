package migrations

var SQLMigration = Sqlx{
	Migrations: []SqlxMigration{
		// Initial Tables
		SqlxFileMigration("create_platforms", "migrations/create_platforms.sql", "migrations/create_platforms.undo.sql"),
		SqlxFileMigration("create_articles", "migrations/create_articles.sql", "migrations/create_articles.undo.sql"),
		SqlxFileMigration("create_categories", "migrations/create_categories.sql", "migrations/create_categories.undo.sql"),
		SqlxFileMigration("create_contacts", "migrations/create_contacts.sql", "migrations/create_contacts.undo.sql"),
		SqlxFileMigration("create_projects", "migrations/create_projects.sql", "migrations/create_projects.undo.sql"),
		SqlxFileMigration("create_tags", "migrations/create_tags.sql", "migrations/create_tags.undo.sql"),
		SqlxFileMigration("create_users", "migrations/create_users.sql", "migrations/create_users.undo.sql"),
		SqlxFileMigration("create_platforms_articles", "migrations/create_platforms_articles.sql", "migrations/create_platforms_articles.undo.sql"),
		SqlxFileMigration("create_platforms_categories", "migrations/create_platforms_categories.sql", "migrations/create_platforms_categories.undo.sql"),
		SqlxFileMigration("create_platforms_projects", "migrations/create_platforms_projects.sql", "migrations/create_platforms_projects.undo.sql"),
		SqlxFileMigration("create_articles_tags", "migrations/create_articles_tags.sql", "migrations/create_articles_tags.undo.sql"),
		SqlxFileMigration("create_projects_tags", "migrations/create_projects_tags.sql", "migrations/create_projects_tags.undo.sql"),

		// Categories
		SqlxFileMigration("insert_categories", "migrations/insert_categories.sql", "migrations/insert_categories.undo.sql"),
	},
}
