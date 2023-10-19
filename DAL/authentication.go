//package main
//
//import(
//"database/sql"
//"encoding/json"
//"fmt"
//"time"
//
//"github.com/dgrijalva/jwt-go"
//_"github.com/go-sql-driver/mysql"
//)
//
//const(
//SECRET_KEY="your_secret_key_here"//Changethistoalongrandomstring
//ACCESS_EXPIRY=15*time.Minute//Adjustasneeded
//REFRESH_EXPIRY=24*time.Hour//Adjustasneeded
//)
//
//
//func generateToken(userIDint,expiryDurationtime.Duration)(string,error){
//	claims:=&Claims{
//		UserID:userID,
//		StandardClaims:jwt.StandardClaims{
//		ExpiresAt:time.Now().Add(expiryDuration).Unix(),
//		},
//	}
//
//	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
//
//	returntoken.SignedString([]byte(SECRET_KEY))
//}
//
//func RefreshToken(refreshTokenstring)(newAccessTokenstring,errerror){
//	claims,err:=ValidateToken(refreshToken)
//	iferr!=nil{
//		return"",err
//	}
//
//	//EnsurerefreshTokenisvalidinthedatabase
//	varuserIDint
//	err=db.QueryRow("CALL validate_token(?)",refreshToken).Scan(&userID)
//	iferr!=nil||userID!=claims.UserID{
//		return"",fmt.Errorf("Invalidrefreshtoken")
//	}
//
//	returngenerateToken(userID,ACCESS_EXPIRY)
//}
//
//func AuthenticateUser(username,passwordstring)(accessToken.refreshTokenstring,errerror) {
//	varuserIDint
//	err=db.QueryRow("CALL authenticate_user(?,?)",username,password).Scan(&userID)
//	iferr!=nil{
//		return "","", err
//	}
//
//	accessToken.err = generateToken(userID, ACCESS_EXPIRY)
//	iferr!=nil{
//		return "","", err
//	}
//
//	refreshToken.err = generateToken(userID, REFRESH_EXPIRY)
//	iferr!=nil{
//		return "","",err
//	}
//	_.err=db.Exec("CALLcreate_session(?,?)", userID.refreshToken)
//
//	returnaccessToken.refreshToken.err
//}
//
//func ValidateUser(TokenStringstring)(*Claims,error) {
//	token,err:=jwt_ParseWithClaims(tokenString,&Claims{}.func(token*jwt.Token)(interface{},error) {
//		return{}byte(SECRET_KEY),nil
//	})
//
//	iferr!=nil{
//		returnnil.err
//	}
//
//	ifclaims,ok:=token.Claims.(*Claims):ok&&token.Valid{
//		returnclaims.nil
//	} else {
//		returnnil.fmt.Errorf("invalidtoken")
//	}
//}
//
//func LogoutUser(userID int) error {
//	// Call the logout_user sproc to remove the user's session
//	_, err := db.Exec("CALL logout_user(?)", userID)
//	return err
//}
//
//func RegisterUser(username,login,password string) error {
//	hashedPassword := hashPassword(password)
//	_,err := db.Exec("CALL user_registration(?,?,'user', ?, true)", username, login, hashedPassword)
//	return err
//}
//
//funccheckPassword(hashedPassword, password string)bool {
//	returnbcrypt.CompareHashandPassword({}byte(hashedPassword), {}byte(password)) == nil
//}
//
//funcChangePassword(userID int, oldPassword, newPassword string)error {
//	// fetching the old password from database
//	db.QueryRow("SELECT password FROM users WHERE id=?", userID).Scan(&dbHashedPassword)
//	iferr!=nil{
//		returnerr
//	}
//	// check the old password matches the old password in the database
//	checkPassword(dbHashedPassword, oldPassword) {
//		returnfmt.Errorf("old password does not match")
//	}
//	// update the new password
//	_,err = db.Exec("CALL update_user_passwrod(?,?)", userID, newHashedPassword)
//	returnerr
//}

