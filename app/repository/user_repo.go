package repository

import (
	"example1/app/model"
	"example1/database"
	"example1/utils/token"
	"time"
	
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type _UserRepository struct {
}

func UserRepository() *_UserRepository {
	return &_UserRepository{}
}

// Login Check
func (h *_UserRepository) CheckUserPassword(condition *model.LoginStudent) (Student model.Student, result *gorm.DB, tokenResult string) {
	name := condition.Name
	student := model.Student{}
	result = database.DB.Where("name = ?", name).First(&student)
	pwdMatch, err := comparePasswords(student.Password, condition.Password)
	if !pwdMatch {
		result.Error = err
		tokenResult = "something error!"
		return student, result, tokenResult
	}

	// Token：若成功搜尋到，產出Token
	tokenResult,err = token.GenerateToken(student.Id)

	if err != nil {
		tokenResult = "something error!"
		return student, result, tokenResult
	}
	
	return student, result, tokenResult
}

// Create User
func (h *_UserRepository) Create(data *model.CreateStudent) (id int, result *gorm.DB) {
	student := model.Student{
		Name:           data.Name,
		Password:       data.Password,
		Student_number: data.Student_number,
		CreatedTime:    time.Now(),
		UpdatedTime:    time.Now()}
	result = database.DB.Create(&student)
	return student.Id, result
}

// score search
func (h *_UserRepository) ScoreSearch(requestData string) (studentInterface []interface{}, studentSearch model.SearchStudent) {
	// 宣告student格式給rows的搜尋結果套用
	student := model.Student{}
	// 將三張資料表join起來，去搜尋是否有id=requestData的人，並拿出指定欄位
	rows, err := database.DB.Model(&student).Select("scores.score,students.name,courses.subject").
		Joins("left join scores on students.id = scores.student_id").
		Joins("left join courses on courses.id = scores.course_id").Where("students.id = ?", requestData).Rows()
	// 如果rows沒找到就不循覽結果直接回傳空interface，如果rows找到就去尋覽結果並傳到新的studentInterface
	if err == nil {
		for rows.Next() {
			database.DB.ScanRows(rows, &studentSearch)
			studentInterface = append(studentInterface, studentSearch)
		}
	}
	// 資料庫最後再關閉
	defer rows.Close()
	return studentInterface, studentSearch
}

// hash 方法
func comparePasswords(hashedPwd string, plainPwd string) (bool, error) {
	byteHash := []byte(hashedPwd)
	byteHash2 := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, byteHash2)
	if err != nil {
		return false, err
	}
	return true, err
}
