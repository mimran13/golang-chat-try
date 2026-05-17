package user

import (
	"database/sql"
	"time"
)

type User struct {
	ID int64
	Username string
	Email string
	CreatedAt time.Time
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
} 

func (r *UserRepository) GetAll() ([]User, error) {
rows, err := r.db.Query("SELECT id, username, email, created_at FROM users")                                                                                                                                             
  if err != nil {                                                                                                                                                                                                          
      return nil, err                                                                                                                                                                                                      
  }                                                                                                                                                                                                                      
  defer rows.Close()                                                                                                                                                                                                       
                    
  var users []User                                                                                                                                                                                                         
  for rows.Next() {                                                                                                                                                                                                      
      var u User   
      if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt); err != nil {
          return nil, err                                                          
      }                  
      users = append(users, u)                                                                                                                                                                                             
  }                           
                                                                                                                                                                                                                           
  return users, rows.Err()    
}