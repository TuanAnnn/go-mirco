package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeOut = time.Second * 3

var db *sql.DB

// New is the function used to create an instance of data package. It return the type
// Model, which embeds all the types we want to be available to our application
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: User{},
	}
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function

type Models struct {
	User User
}

// User is the structure with holds one user from the database
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:created_at`
	UpdatedAt time.Time `json:updated_at`
}

// get all returns a slice of all user, sorted by last name
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	from users order by last_name`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}
	return users, nil
}

// getByEmail returns one user by email

func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	// In ra giá trị email được truyền vào để kiểm tra
	log.Printf("Executing GetByEmail with email: %s", email)

	// Câu lệnh SQL với điều kiện lọc email
	query := `SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at 
	          FROM users
	          WHERE email = $1`

	var user User

	// Thực hiện truy vấn với giá trị email
	row := db.QueryRowContext(ctx, query, email)

	// Quét dữ liệu từ kết quả trả về và ghi log trước khi quét
	log.Printf("Query executed, scanning result for email: %s", email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	// Xử lý lỗi nếu có
	if err != nil {
		if err == sql.ErrNoRows {
			// Ghi log nếu không tìm thấy người dùng
			log.Printf("No user found with email: %s", email)
			return nil, errors.New("no user found with that email")
		}
		// Ghi log lỗi khác
		log.Printf("Error scanning user with email: %s, error: %v", email, err)
		return nil, err
	}

	// Ghi log nếu tìm thấy người dùng
	log.Printf("User found with email: %s, ID: %d", user.Email, user.ID)

	// Trả về người dùng nếu tìm thấy
	return &user, nil
}

// get one user by user by id

func (u *User) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	// Câu lệnh SQL có điều kiện lọc theo id
	query := `SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at 
	          FROM users 
	          WHERE id = $1`

	var user User

	// Thực hiện truy vấn với tham số id
	row := db.QueryRowContext(ctx, query, id)

	// Quét dữ liệu từ kết quả truy vấn
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	// Xử lý lỗi nếu có
	if err != nil {
		if err == sql.ErrNoRows {
			// Trả về nil và lỗi nếu không tìm thấy người dùng
			return nil, errors.New("no user found with that ID")
		}
		return nil, err
	}

	// Trả về người dùng nếu tìm thấy
	return &user, nil
}

// update updates one user in the database, using the interformation
// stored in the receiver u
func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `update users set
	email = $1,
	first_name = $2,
	last_name = $3,
	user_active = $4,
	updated_at = $5,
	where id = $6
	`

	_, err := db.ExecContext(ctx, stmt,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Active,
		time.Now(),
		u.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// Delete deletes one user from the database, by user.ID

func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByID deletes one user from the database, by ID
func (u *User) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Insert(user User) (int, error) {
	// Tạo một context với timeout
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	// Hash mật khẩu người dùng với bcrypt và log lỗi nếu có
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	// Lấy thời gian hiện tại để sử dụng cho cả created_at và updated_at
	now := time.Now()

	log.Printf("Inserting user: %s, %s, %s", user.Email, user.FirstName, user.LastName)

	// Câu lệnh SQL chèn người dùng mới vào cơ sở dữ liệu
	stmt := `insert into public.users (email, first_name, last_name, password, user_active, created_at, updated_at)
			 values ($1, $2, $3, $4, $5, $6, $7) returning id`

	var newId int

	// Thực hiện câu lệnh chèn với các tham số và lấy id mới
	err = db.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.Active,
		now,
		now,
	).Scan(&newId)

	if err != nil {
		return 0, err
	}

	return newId, nil
}

// Reset password is the method we will use to change a user's password

func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `update users set password = $1 where id = $2`

	_, err = db.ExecContext(ctx, stmt, hashedPassword, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
