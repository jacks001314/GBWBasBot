package security

import (
	"os/user"
	"sync"

	"common/proto/hadoop/hadoop_common"
)

/** a (very) basic UserGroupInformation implementation for storing user data/tokens,
  This implementation is currently *not* thread-safe
*/

type UserGroupInformation struct {
	rwMutex    sync.RWMutex
	userInfo   *hadoop_common.UserInformationProto
	userTokens map[string]*hadoop_common.TokenProto
}

var once sync.Once
var currentUserGroupInformation *UserGroupInformation
var maxTokens = 16

func CreateCurrentUserInfoProto() (*hadoop_common.UserInformationProto, error) {
	// Figure the current user-name
	var username string
	if currentUser, err := user.Current(); err != nil {

		return nil, err
	} else {
		username = currentUser.Username
	}

	return &hadoop_common.UserInformationProto{EffectiveUser: nil, RealUser: &username}, nil
}

func Allocate(userInfo *hadoop_common.UserInformationProto, userTokens map[string]*hadoop_common.TokenProto) *UserGroupInformation {
	ugi := new(UserGroupInformation)

	if userInfo != nil {
		ugi.userInfo = userInfo
	} else {
		currentUserInfo, _ := CreateCurrentUserInfoProto()
		ugi.userInfo = currentUserInfo
	}

	if userTokens != nil {
		ugi.userTokens = userTokens
	} else {
		ugi.userTokens = make(map[string]*hadoop_common.TokenProto) //empty, with room for maxTokens tokens.
	}

	return ugi
}

func initializeCurrentUser() {
	once.Do(func() {
		currentUserGroupInformation = Allocate(nil, nil)
	})
}

func (ugi *UserGroupInformation) GetUserInformation() *hadoop_common.UserInformationProto {
	return ugi.userInfo
}

func (ugi *UserGroupInformation) GetUserTokens() map[string]*hadoop_common.TokenProto {
	return ugi.userTokens
}

func (ugi *UserGroupInformation) AddUserTokenWithAlias(alias string, token *hadoop_common.TokenProto) {
	if token == nil {

		return
	}

	if length := len(ugi.userTokens); length < maxTokens {
		ugi.userTokens[alias] = token
	} else {

	}
}

func (ugi *UserGroupInformation) AddUserToken(token *hadoop_common.TokenProto) {
	if token == nil {

		return
	}

	ugi.AddUserTokenWithAlias(token.GetService(), token)
}

func GetCurrentUser() *UserGroupInformation {
	initializeCurrentUser()

	return currentUserGroupInformation
}
