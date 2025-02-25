# 🚀 Go Rajaongkir Destination Cache

This Go service is responsible for fetching destinations from a local database. If the requested destination is not found, it forwards the request to the RajaOngkir API, stores the response in the database, and returns it to the requester.

## 📑 Features
- Fetch destinations from the local PostgreSQL database
- Query RajaOngkir API when a destination is not found
- Store new destinations in the database for future use
- Uses GORM as ORM for database interactions
- Optimized search by zip code

---

## 🛠️ Setup & Installation

### **Prerequisites**
- Go 1.20+ installed
- PostgreSQL installed and running
- `.env` file with required configurations

### **1️⃣ Clone the repository**
```sh
git clone https://github.com/yourusername/go-destination-service.git
cd go-destination-service
```


### **2️⃣ Install Depedencies**
```sh
go mod tidy
```

### **3️⃣ Run Database Migrations**
```sh
go run main.go migrate
```


### **4️⃣ Start the server**
```sh
go run main.go
```

or 

```sh
air
```