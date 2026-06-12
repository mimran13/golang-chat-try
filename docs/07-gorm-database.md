# 07 - GORM Database ORM

GORM is Go's most popular ORM (Object-Relational Mapping). If you've used TypeORM or Sequelize, you'll find GORM familiar.

## Why GORM?

| Feature | GORM | TypeORM | Sequelize |
|---------|------|---------|-----------|
| Language | Go | TypeScript | JavaScript |
| Pattern | Active Record + Data Mapper | Data Mapper | Active Record |
| Migrations | Auto-migrate + manual | CLI migrations | CLI migrations |
| Relations | Has One/Many, Belongs To, M2M | Same | Same |
| Hooks | BeforeCreate, AfterCreate, etc. | Subscribers | Hooks |

---

## Connection Setup

See `internal/database/database.go`:

**TypeORM (NestJS):**
```typescript
@Module({
  imports: [
    TypeOrmModule.forRoot({
      type: 'mysql',
      host: 'localhost',
      port: 3306,
      username: 'root',
      password: 'password',
      database: 'task_manager',
      entities: [User, Task],
      synchronize: true,
    }),
  ],
})
export class AppModule {}
```

**GORM (Go):**
```go
import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // Configure connection pool
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    return db, nil
}
```

---

## Model Definition

### Basic Model

**TypeORM Entity:**
```typescript
@Entity()
export class User {
  @PrimaryGeneratedColumn()
  id: number;

  @Column({ unique: true })
  username: string;

  @Column({ unique: true })
  email: string;

  @Column()
  passwordHash: string;

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;

  @DeleteDateColumn()
  deletedAt: Date;
}
```

**GORM Model:**
```go
type User struct {
    gorm.Model               // Includes ID, CreatedAt, UpdatedAt, DeletedAt
    Username     string      `gorm:"uniqueIndex;size:50;not null"`
    Email        string      `gorm:"uniqueIndex;size:100;not null"`
    PasswordHash string      `gorm:"size:255;not null"`
    Tasks        []Task      `gorm:"foreignKey:UserID"`  // Has many
}

// gorm.Model provides:
// ID        uint           `gorm:"primaryKey"`
// CreatedAt time.Time
// UpdatedAt time.Time
// DeletedAt gorm.DeletedAt `gorm:"index"`  // For soft delete
```

### GORM Tags

```go
type Task struct {
    ID          uint           `gorm:"primaryKey"`
    UserID      uint           `gorm:"not null;index"`
    Title       string         `gorm:"size:200;not null"`
    Description string         `gorm:"type:text"`
    Status      TaskStatus     `gorm:"type:varchar(20);default:pending"`
    Priority    TaskPriority   `gorm:"type:varchar(20);default:medium"`
    DueDate     *time.Time     `gorm:"index"`                  // Nullable
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
    User        User           `gorm:"foreignKey:UserID"`      // Belongs to
}
```

Common GORM tags:
| Tag | Description |
|-----|-------------|
| `primaryKey` | Primary key |
| `uniqueIndex` | Unique index |
| `index` | Regular index |
| `not null` | NOT NULL constraint |
| `size:100` | VARCHAR(100) |
| `type:text` | Specific SQL type |
| `default:value` | Default value |
| `foreignKey:Field` | Foreign key |
| `-` | Ignore field |

---

## Auto Migration

**TypeORM (synchronize):**
```typescript
TypeOrmModule.forRoot({
  synchronize: true,  // Auto-creates tables
})
```

**GORM:**
```go
// Auto-migrate creates/updates tables
err := db.AutoMigrate(&User{}, &Task{})
if err != nil {
    log.Fatal("Migration failed:", err)
}
```

**What AutoMigrate does:**
- Creates tables if they don't exist
- Adds missing columns
- Adds missing indexes
- Does NOT delete columns or change types (safe)

---

## CRUD Operations

### Create

**TypeORM:**
```typescript
const user = userRepository.create({ username: 'john', email: 'john@example.com' });
await userRepository.save(user);
```

**GORM:**
```go
user := &User{
    Username:     "john",
    Email:        "john@example.com",
    PasswordHash: hashedPassword,
}
result := db.Create(user)  // user.ID is set after this

if result.Error != nil {
    return result.Error
}
fmt.Println("Rows affected:", result.RowsAffected)
fmt.Println("New ID:", user.ID)
```

### Read

**TypeORM:**
```typescript
// Find by ID
const user = await userRepository.findOne({ where: { id: 1 } });

// Find all with conditions
const users = await userRepository.find({
  where: { status: 'active' },
  order: { createdAt: 'DESC' },
  skip: 0,
  take: 10,
});
```

**GORM:**
```go
// Find by ID
var user User
result := db.First(&user, 1)  // Find by primary key
// OR
result := db.First(&user, "id = ?", 1)

if errors.Is(result.Error, gorm.ErrRecordNotFound) {
    // Not found
}

// Find all with conditions
var users []User
db.Where("status = ?", "active").
    Order("created_at DESC").
    Offset(0).
    Limit(10).
    Find(&users)

// Find one with conditions
var user User
db.Where("email = ?", "john@example.com").First(&user)
```

### Update

**TypeORM:**
```typescript
await userRepository.update(1, { username: 'johnny' });
// OR
user.username = 'johnny';
await userRepository.save(user);
```

**GORM:**
```go
// Update single field
db.Model(&user).Update("username", "johnny")

// Update multiple fields
db.Model(&user).Updates(User{Username: "johnny", Email: "new@example.com"})

// Update with map (includes zero values)
db.Model(&user).Updates(map[string]interface{}{
    "username": "johnny",
    "age":      0,  // This WILL be updated
})

// Save entire struct
user.Username = "johnny"
db.Save(&user)
```

**Important:** `Updates` with struct ignores zero values! Use map for explicit zeros.

### Delete

**TypeORM:**
```typescript
await userRepository.softDelete(1);  // Soft delete
await userRepository.delete(1);      // Hard delete
```

**GORM:**
```go
// Soft delete (if model has DeletedAt field)
db.Delete(&user)  // Sets deleted_at, doesn't actually delete

// Hard delete (permanent)
db.Unscoped().Delete(&user)

// Delete by ID
db.Delete(&User{}, 1)
db.Delete(&User{}, []int{1, 2, 3})  // Multiple
```

---

## Querying

### Where Conditions

```go
// Simple conditions
db.Where("name = ?", "john").Find(&users)
db.Where("name <> ?", "john").Find(&users)
db.Where("name IN ?", []string{"john", "jane"}).Find(&users)
db.Where("name LIKE ?", "%john%").Find(&users)
db.Where("age > ?", 18).Find(&users)
db.Where("created_at BETWEEN ? AND ?", start, end).Find(&users)

// Struct conditions (ignores zero values!)
db.Where(&User{Username: "john", Age: 0}).Find(&users)
// Only queries: WHERE username = 'john'

// Map conditions (includes zero values)
db.Where(map[string]interface{}{"username": "john", "age": 0}).Find(&users)
// Queries: WHERE username = 'john' AND age = 0

// Multiple conditions
db.Where("status = ?", "active").
    Where("age > ?", 18).
    Find(&users)

// OR conditions
db.Where("role = ?", "admin").
    Or("role = ?", "super_admin").
    Find(&users)

// NOT conditions
db.Not("name = ?", "john").Find(&users)
```

### Select Specific Fields

```go
// Select specific columns
db.Select("id", "username").Find(&users)

// Select with expression
db.Select("id", "CONCAT(first_name, ' ', last_name) as full_name").Find(&users)
```

### Ordering and Pagination

```go
// Order
db.Order("created_at DESC").Find(&tasks)
db.Order("priority DESC, created_at ASC").Find(&tasks)

// Pagination
page := 1
pageSize := 10
offset := (page - 1) * pageSize

db.Offset(offset).Limit(pageSize).Find(&tasks)

// Count total
var total int64
db.Model(&Task{}).Where("user_id = ?", userID).Count(&total)
```

### Joins and Preloading

**TypeORM:**
```typescript
const tasks = await taskRepository.find({
  relations: ['user'],  // Eager load
});
```

**GORM:**
```go
// Preload (separate query)
var tasks []Task
db.Preload("User").Find(&tasks)

// Preload with conditions
db.Preload("Tasks", "status = ?", "pending").Find(&users)

// Nested preload
db.Preload("Tasks.Comments").Find(&users)

// Joins (single query)
db.Joins("User").Find(&tasks)

// Raw joins
db.Joins("LEFT JOIN users ON users.id = tasks.user_id").
    Where("users.status = ?", "active").
    Find(&tasks)
```

---

## Hooks (Lifecycle Callbacks)

**TypeORM:**
```typescript
@Entity()
export class User {
  @BeforeInsert()
  hashPassword() {
    this.password = bcrypt.hashSync(this.password);
  }
}
```

**GORM:**
```go
// Available hooks:
// BeforeSave, BeforeCreate, BeforeUpdate, BeforeDelete
// AfterSave, AfterCreate, AfterUpdate, AfterDelete
// AfterFind

func (u *User) BeforeCreate(tx *gorm.DB) error {
    // Hash password before saving
    if u.PasswordHash == "" {
        return errors.New("password required")
    }
    return nil
}

func (t *Task) BeforeUpdate(tx *gorm.DB) error {
    // Validate before update
    if t.Title == "" {
        return errors.New("title cannot be empty")
    }
    return nil
}

func (u *User) AfterFind(tx *gorm.DB) error {
    // Transform after fetching
    u.Email = strings.ToLower(u.Email)
    return nil
}
```

---

## Transactions

**TypeORM:**
```typescript
await dataSource.transaction(async (manager) => {
  await manager.save(user);
  await manager.save(task);
});
```

**GORM:**
```go
// Automatic transaction
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err  // Rollback
    }
    if err := tx.Create(&task).Error; err != nil {
        return err  // Rollback
    }
    return nil  // Commit
})

// Manual transaction
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Create(&task).Error; err != nil {
    tx.Rollback()
    return err
}

tx.Commit()
```

---

## Context Support

GORM supports context for timeouts and cancellation:

```go
// With context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

var users []User
db.WithContext(ctx).Find(&users)

// In handlers (using Gin's context)
func (h *Handler) GetUsers(c *gin.Context) {
    var users []User
    err := h.db.WithContext(c.Request.Context()).Find(&users).Error
    // Request cancellation will cancel the query
}
```

---

## Our Repository Pattern

See `internal/repository/`:

```go
// Interface definition
type UserRepository interface {
    FindByID(ctx context.Context, id uint) (*models.User, error)
    FindByEmail(ctx context.Context, email string) (*models.User, error)
    Create(ctx context.Context, user *models.User) error
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id uint) error
}

// Implementation
type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
    var user models.User
    err := r.db.WithContext(ctx).First(&user, id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil  // Not found, not an error
    }
    return &user, err
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}
```

---

## Exercises

1. Look at `internal/models/task.go` - what GORM tags are used?
2. Look at `internal/repository/task_repository.go` - how is pagination implemented?
3. Add a new field `CompletedAt *time.Time` to the Task model
4. Write a query to find all overdue tasks (due_date < now AND status != completed)
