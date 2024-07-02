package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

// var db *sql.DB

type Employees struct {
	Emp_no     int    `json:"emp_no"`
	Birth_date string `json:"birth_date"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Gender     string `json:"gender"`
	Hire_date  string `json:"hire_date"`
}

type UpdateEmployees struct {
	Birth_date string `json:"birth_date"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}

// var result string

func main() {

	/// connect to the database ///

	dsn := "root:12345@tcp(127.0.0.1:3306)/employees_mod"

	db, err := sql.Open("mysql", dsn) ///
	if err != nil {
		log.Fatal("Failed to open database", err)
	}
	fmt.Println("Opened Successfuly")

	defer db.Close()

	server := &APIServer{

		addr: ":8000",
		db:   db,
	}

	router := gin.Default()

	router.GET("/employees", server.GetEmployees)

	router.POST("/employees", server.AddEmployees)

	router.DELETE("/employees/:emp_no", server.DeleteEmployee)

	router.PATCH("/employees/:emp_no", server.UpdateEmployees)

	log.Printf("Starting server on %s", server.addr)
	if err := router.Run(server.addr); err != nil {
		log.Fatal("Failed", err)
	}

}

func (s *APIServer) GetEmployees(c *gin.Context) {

	rows, err := s.db.Query("Select * from t_employees")
	if err != nil {
		log.Println("error in querying employees ", err)

		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

		return
	}
	defer rows.Close()

	var employees []Employees
	for rows.Next() {

		var emp Employees
		err := rows.Scan(&emp.Emp_no, &emp.Birth_date, &emp.First_name, &emp.Last_name, &emp.Gender, &emp.Hire_date)
		if err != nil {
			log.Println("Error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process employee data"})

			return
		}
		employees = append(employees, emp)

	}
	c.JSON(http.StatusOK, employees)

}

func (s *APIServer) AddEmployees(c *gin.Context) {

	var NewEmployee Employees

	if err := c.ShouldBindJSON(&NewEmployee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	}

	result, err := s.db.Exec("insert into t_employees (emp_no, birth_date, first_name, last_name, gender, hire_date) values ( ?, ?, ?, ?, ?, ?)", NewEmployee.Emp_no, NewEmployee.Birth_date, NewEmployee.First_name, NewEmployee.Last_name, NewEmployee.Gender, NewEmployee.Hire_date)
	if err != nil {
		log.Println("Error in inserting employee", err)

		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to add employee"})

		return
	}
	c.JSON(http.StatusOK, NewEmployee)

	fmt.Println(result)

}

func (s *APIServer) DeleteEmployee(c *gin.Context) {

	empid := c.Param("emp_no")

	id, err := strconv.Atoi(empid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid ID", "ID": empid})
	}

	result, err := s.db.Exec("Delete from t_employees where emp_no = ?", id)
	if err != nil {

		log.Println("Error deleting employee:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete employee"})
		return

	}

	rowsaffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal()
	}

	c.JSON(http.StatusOK, gin.H{"Message": fmt.Sprintf("Deleted Successfuly ID %d", id)})
	c.JSON(http.StatusOK, gin.H{"Message": fmt.Sprintf("Numeber of rows affected: %d", rowsaffected)})
	log.Println("Deleted employee successfuly", result)

}

func (s *APIServer) UpdateEmployees(c *gin.Context) {

	empid := c.Param("emp_no")

	id, err := strconv.Atoi(empid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid ID", "ID": empid})
	}

	var UpEmp UpdateEmployees
	if err := c.ShouldBindJSON(&UpEmp); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to update "})

	}
	fmt.Print(UpEmp.Last_name)
	result, err := s.db.Exec("update t_employees set birth_date = ? ,  first_name = ? , last_name = ?  where emp_no = ? ", UpEmp.Birth_date, UpEmp.First_name, UpEmp.Last_name, id)

	if err != nil {
		log.Println("Error in updating employee", err)

		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to add employee"})

		return
	}

	rowsaffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Number or rows affected: %d\n", rowsaffected)

	c.JSON(http.StatusOK, gin.H{"Message": "Updated Successfuly"})
	c.JSON(http.StatusOK, gin.H{"Message": fmt.Sprintf("Number of rows affected: %d", rowsaffected)})
	// log.Println("updated employee", result)

}
