# Order Service

Order Service เป็นส่วนหนึ่งของระบบ E-Commerce Microservice ที่จัดการเกี่ยวกับการสั่งซื้อสินค้า โดยใช้ Clean Architecture ในการพัฒนา

## โครงสร้างโปรเจค (Project Structure)

```
order-service/
├── internal/
│   ├── domain/          # Business entities และ interfaces
│   ├── repository/      # Data access layer implementations
│   ├── usecase/        # Business logic implementations
│   └── delivery/       # HTTP handlers
├── main.go             # Entry point
└── README.md
```

## Clean Architecture

โปรเจคนี้ใช้ Clean Architecture ในการพัฒนา แบ่งเป็น 4 layer หลัก:

### 1. Domain Layer (internal/domain/)
- เก็บ Business entities และ interfaces
- กำหนด contracts ระหว่าง layers
- ไม่มีการพึ่งพา external dependencies

### 2. Repository Layer (internal/repository/)
- จัดการการเข้าถึงฐานข้อมูล
- implement interfaces จาก Domain Layer
- รับผิดชอบการแปลงข้อมูลระหว่าง entities และ database model

### 3. Use Case Layer (internal/usecase/)
- จัดการ Business logic
- ใช้ interfaces จาก Domain Layer
- ไม่รู้จัก implementation details ของ Repository

### 4. Delivery Layer (internal/delivery/)
- จัดการ HTTP handlers
- แปลง HTTP requests/responses
- เรียกใช้ Use Cases ผ่าน interfaces

## Flow การทำงาน

ตัวอย่างการสร้าง Order ใหม่:

1. **HTTP Request**
```http
POST /api/v1/orders
{
    "user_id": "123",
    "product_id": "456",
    "quantity": 2,
    "total_price": 1000
}
```

2. **Delivery Layer** (HTTP Handler)
```go
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    // 1. รับ request และแปลงเป็น domain.Order
    var order domain.Order
    json.NewDecoder(r.Body).Decode(&order)

    // 2. ส่งต่อไปให้ Use Case
    h.orderUseCase.CreateOrder(r.Context(), &order)
}
```

3. **Use Case Layer**
```go
func (u *orderUseCase) CreateOrder(ctx context.Context, order *domain.Order) error {
    // 3. จัดการ business logic
    order.CreatedAt = time.Now()
    order.UpdatedAt = time.Now()
    order.Status = "pending"

    // 4. ส่งต่อไปให้ Repository
    return u.orderRepo.Create(ctx, order)
}
```

4. **Repository Layer**
```go
func (r *mongoOrderRepository) Create(ctx context.Context, order *domain.Order) error {
    // 5. บันทึกลงฐานข้อมูล
    _, err := r.collection.InsertOne(ctx, order)
    return err
}
```

5. **Response**
```json
{
    "id": "507f1f77bcf86cd799439011",
    "user_id": "123",
    "product_id": "456",
    "quantity": 2,
    "total_price": 1000,
    "status": "pending",
    "created_at": "2024-03-06T12:00:00Z",
    "updated_at": "2024-03-06T12:00:00Z"
}
```

## API Endpoints

- `POST /api/v1/orders` - สร้าง order ใหม่
- `GET /api/v1/orders` - ดึงรายการ orders ทั้งหมด
- `GET /api/v1/orders/{id}` - ดึงข้อมูล order ตาม ID
- `PUT /api/v1/orders/{id}` - อัพเดท order
- `DELETE /api/v1/orders/{id}` - ลบ order
- `GET /health` - Health check endpoint

## การติดตั้งและรัน

1. ติดตั้ง dependencies:
```bash
go mod download
```

2. ตั้งค่า environment variables:
```bash
export MONGODB_URI="mongodb://localhost:27017"
export PORT="8083"
```

3. รัน service:
```bash
go run main.go
```

## ข้อดีของ Clean Architecture

1. **Separation of Concerns**
   - แต่ละ layer มีหน้าที่ชัดเจน
   - ลดความซับซ้อนของโค้ด
   - ง่ายต่อการดูแลรักษา

2. **Dependency Rule**
   - Layer ชั้นในไม่รู้จัก Layer ชั้นนอก
   - Domain Layer เป็นศูนย์กลาง
   - ใช้ Interface ในการสื่อสารระหว่าง Layer

3. **Testability**
   - สามารถ Mock แต่ละ Layer ได้ง่าย
   - Test แยกแต่ละ Layer ได้อิสระ
   - ครอบคลุมการทดสอบทั้ง Unit Test และ Integration Test

4. **Flexibility**
   - สามารถเปลี่ยน implementation ได้ง่าย
   - รองรับการขยายระบบในอนาคต
   - ลดผลกระทบเมื่อต้องเปลี่ยนแปลง external dependencies
