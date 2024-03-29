package db

type Contact struct {
	Model
	Name       string `json:"name" db:"name" binding:"required"`
	Title      string `json:"title" db:"title"`
	Email      string `json:"email" db:"email"`
	Phone      string `json:"phone" db:"phone"`
	Phone2     string `json:"phone2" db:"phone2"`
	Address    string `json:"address" db:"address"`
	Notes      string `json:"notes" db:"notes"`
	Source     string `json:"source" db:"source"`
	Privacy    string `json:"privacy" db:"privacy"`
	PlatformId int64  `json:"platformId" db:"platform_id"`
}

func (db *Database) CountContacts() (int, error) {
	var count int
	err := db.querier.Get(&count, "SELECT COUNT(*) AS count FROM contacts WHERE deleted_at IS NULL")
	return count, err
}

func (db *Database) GetContactsForPlatform(platformId int64) ([]Contact, error) {
	contacts := []Contact{}
	err := db.querier.Select(&contacts, "SELECT * FROM contacts WHERE platform_id = ? AND deleted_at IS NULL", platformId)
	return contacts, err
}

func (db *Database) EditContact(contact Contact) error {
	_, err := db.querier.NamedExec(`
		UPDATE contacts 
		SET 
			name = :name, 
			title = :title, 
			email = :email, 
			phone = :phone, 
			phone2 = :phone2, 
			address = :address, 
			notes = :notes, 
			source = :source, 
			privacy = :privacy 
		WHERE id = :id AND platform_id = :platform_id`, contact)
	return err
}

func (db *Database) InsertContact(contact Contact) error {
	_, err := db.querier.NamedExec(`
		INSERT INTO contacts 
			(name, title, email, phone, phone2, address, notes, source, privacy, platform_id) 
		VALUES 
			(:name, :title, :email, :phone, :phone2, :address, :notes, :source, :privacy, :platform_id)`, contact)
	return err
}

func (db *Database) DeleteContact(id, platformId int64) error {
	_, err := db.querier.Exec("UPDATE contacts SET deleted_at = CURRENT_TIMESTAMP() WHERE id = ? AND platform_id = ?", id, platformId)
	return err
}
