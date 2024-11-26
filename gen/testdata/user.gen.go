package testdata

import (
	"github.com/DaHuangQwQ/gsql"

	"database/sql"
)

const (
	UserName     = "Name"
	UserAge      = "Age"
	UserNickName = "NickName"
	UserPicture  = "Picture"
)

func UserNameEq(val string) gsql.Predicate {
	return gsql.C("Name").Eq(val)
}

func UserAgeEq(val *int) gsql.Predicate {
	return gsql.C("Age").Eq(val)
}

func UserNickNameEq(val *sql.NullString) gsql.Predicate {
	return gsql.C("NickName").Eq(val)
}

func UserPictureEq(val []byte) gsql.Predicate {
	return gsql.C("Picture").Eq(val)
}

const (
	UserDetailAddress = "Address"
)

func UserDetailAddressEq(val string) gsql.Predicate {
	return gsql.C("Address").Eq(val)
}
