package lr

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Organizations struct {
	Id                       string              `bson:"_id,omitempty" json:"Id"`
	IsActive                 bool                `bson:"IsActive" json:"IsActive"`
	IsDeleted                bool                `bson:"IsDeleted" json:"-"`
	Name                     string              `bson:"Name" json:"Name"`
	Display                  *Display            `bson:"Display,omitempty" json:"Display,omitempty"`
	Metadata                 map[string]string   `bson:"Metadata,omitempty" json:"Metadata"`
	IsB2B                    bool                `bson:"IsB2B,omitempty" json:"IsB2B"`
	AppId                    int                 `bson:"AppId,omitempty" json:"AppId"`
	Domains                  []Domain            `bson:"Domains,omitempty" json:"Domains,omitempty"`
	IsAuthRestrictedToDomain bool                `bson:"IsAuthRestrictedToDomain,omitempty" json:"IsAuthRestrictedToDomain"`
	Policies                 *Policies           `bson:"Policies,omitempty" json:"Policies,omitempty"`
	Connections              []ConnectionSetting `bson:"Connections,omitempty" json:"Connections,omitempty"`
	CreatedDate              *TimeStamp          `bson:"CreatedDate" json:"CreatedDate"`
	ModifiedDate             *TimeStamp          `bson:"ModifiedDate,omitempty" json:"ModifiedDate,omitempty"`
	CreatedBy                *string             `bson:"CreatedBy,omitempty" json:"CreatedBy,omitempty"`
}

type OrganizationResponse struct {
	Id          string    `json:"Id"`
	RoleId      string    `json:"RoleId"`
	OrgId       string    `json:"OrgId"`
	Uid         string    `json:"Uid"`
	Name        string    `json:"Name"`
	CreatedDate time.Time `json:"CreatedDate"`
}

type Display struct {
	Name    string `bson:"Name" json:"Name"`
	LogoURL string `bson:"LogoURL" json:"LogoURL"`
}

type Domain struct {
	Id                   primitive.ObjectID `bson:"_id" json:"Id"`
	DomainName           string             `bson:"DomainName" json:"DomainName"`
	IsVerified           bool               `bson:"IsVerified" json:"IsVerified"`
	VerificationToken    string             `bson:"VerificationToken,omitempty" json:"VerificationToken,omitempty"`
	VerificationStrategy string             `bson:"VerificationStrategy,omitempty" json:"VerificationStrategy,omitempty"`
}

type Policies struct {
	PasswordPolicy *PasswordPolicies `bson:"PasswordPolicy,omitempty" json:"PasswordPolicy,omitempty"`
	MFAPolicy      *MFAPolicies      `bson:"MFAPolicy,omitempty" json:"MFAPolicy,omitempty"`
	SessionPolicy  *SessionPolicies  `bson:"SessionPolicy,omitempty" json:"SessionPolicy,omitempty"`
	MemberPolicy   *MemberPolicies   `bson:"MemberPolicy,omitempty" json:"MemberPolicy,omitempty"`
	JITPolicy      *JITPolicies      `bson:"JITPolicy,omitempty" json:"JITPolicy,omitempty"`
}

type PasswordPolicies struct {
	MinLength          int  `bson:"MinLength" json:"MinLength"`
	MaxLength          int  `bson:"MaxLength" json:"MaxLength"`
	RequireUppercase   bool `bson:"RequireUppercase" json:"RequireUppercase"`
	RequireLowercase   bool `bson:"RequireLowercase" json:"RequireLowercase"`
	RequireNumber      bool `bson:"RequireNumber" json:"RequireNumber"`
	RequireSpecialChar bool `bson:"RequireSpecialChar" json:"RequireSpecialChar"`
	ExpiryDays         int  `bson:"ExpiryDays" json:"ExpiryDays"`
}

type MFAPolicies struct {
	EnforcementMode string `bson:"EnforcementMode,omitempty" json:"EnforcementMode,omitempty"`
}

type SessionPolicies struct {
	AccessTokenTTL  int `bson:"AccessTokenTTL" json:"AccessTokenTTL"`
	RefreshTokenTTL int `bson:"RefreshTokenTTL" json:"RefreshTokenTTL"`
}

type MemberPolicies struct {
	DefaultMemberRole primitive.ObjectID `bson:"DefaultMemberRole,omitempty" json:"DefaultMemberRole"`
}

type JITPolicies struct {
	Enabled bool `bson:"Enabled" json:"Enabled"`
}

type ConnectionSetting struct {
	Id             primitive.ObjectID `bson:"_id" json:"Id"`
	IsActive       bool               `bson:"IsActive" json:"IsActive"`
	IsDeleted      bool               `bson:"IsDeleted" json:"-"`
	Name           string             `bson:"Name" json:"Name"`
	ConnectionType string             `bson:"ConnectionType" json:"ConnectionType"`
	Domain         string             `bson:"Domain" json:"Domain"`
	Attributes     AttributeMapping   `bson:"Attributes" json:"Attributes"`
	GroupRoles     []GroupRole        `bson:"GroupRoles,omitempty" json:"GroupRoles"`
	CreatedDate    *TimeStamp         `bson:"CreatedDate" json:"CreatedDate"`
	ModifiedDate   *TimeStamp         `bson:"ModifiedDate,omitempty" json:"ModifiedDate,omitempty"`

	//SamlConfig
	EntityId         string          `bson:"EntityId,omitempty" json:"EntityId,omitempty"`
	MetadataUrl      string          `bson:"MetadataUrl,omitempty" json:"MetadataUrl,omitempty"`
	ACSEndpoint      string          `bson:"ACSEndpoint,omitempty" json:"ACSEndpoint,omitempty"`
	SPCertificate    *SPCertificate  `bson:"SPCertificate,omitempty" json:"SPCertificate,omitempty"`
	IDPMetadataUrl   string          `bson:"IDPMetadataUrl,omitempty" json:"IDPMetadataUrl,omitempty"`
	IDPEntityId      string          `bson:"IDPEntityId,omitempty" json:"IDPEntityId,omitempty"`
	IDPLoginUrl      string          `bson:"IDPLoginUrl,omitempty" json:"IDPLoginUrl,omitempty"`
	IDPLoginBinding  string          `bson:"IDPLoginBinding,omitempty" json:"IDPLoginBinding,omitempty"`
	IDPLogoutUrl     string          `bson:"IDPLogoutUrl,omitempty" json:"IDPLogoutUrl,omitempty"`
	IDPLogoutBinding string          `bson:"IDPLogoutBinding,omitempty" json:"IDPLogoutBinding,omitempty"`
	IDPCertificate   *IDPCertificate `bson:"IDPCertificate,omitempty" json:"IDPCertificate,omitempty"`
	IsIDPInitiated   bool            `bson:"IsIDPInitiated,omitempty" json:"IsIDPInitiated,omitempty"`

	//OpenIdConfig
	Issuer           string   `bson:"Issuer,omitempty" json:"Issuer,omitempty"`
	AuthorizationUrl string   `bson:"AuthorizationUrl,omitempty" json:"AuthorizationUrl,omitempty"`
	TokenUrl         string   `bson:"TokenUrl,omitempty" json:"TokenUrl,omitempty"`
	UserInfoUrl      string   `bson:"UserInfoUrl,omitempty" json:"UserInfoUrl,omitempty"`
	ClientId         string   `bson:"ClientId,omitempty" json:"ClientId,omitempty"`
	ClientSecret     string   `bson:"ClientSecret,omitempty" json:"ClientSecret,omitempty"`
	RedirectURI      string   `bson:"RedirectURI,omitempty" json:"RedirectURI,omitempty"`
	Scopes           []string `bson:"Scopes,omitempty" json:"Scopes,omitempty"`
	TokenAuthMethod  string   `bson:"TokenAuthMethod,omitempty" json:"TokenAuthMethod,omitempty"`
}

type AttributeMapping struct {
	ID            string            `bson:"ID,omitempty" json:"ID"`
	Email         string            `bson:"Email,omitempty" json:"Email"`
	FirstName     string            `bson:"FirstName,omitempty" json:"FirstName"`
	LastName      string            `bson:"LastName,omitempty" json:"LastName"`
	Group         string            `bson:"Group,omitempty" json:"Group"`
	CustomMapping map[string]string `bson:"CustomMapping,omitempty" json:"CustomMapping"`
}

type GroupRole struct {
	Id      primitive.ObjectID `bson:"_id" json:"Id"`
	Name    string             `bson:"Name" json:"Name"`       //Unique
	GroupId string             `bson:"GroupId" json:"GroupId"` //Unique
	RoleId  primitive.ObjectID `bson:"RoleId" json:"RoleId"`
}

type SPCertificate struct {
	Key         string `bson:"Key" json:"-"`
	Certificate string `bson:"Certificate,omitempty" json:"Certificate,omitempty"`
}

type IDPCertificate struct {
	Certificate string    `bson:"Certificate,omitempty" json:"Certificate,omitempty"`
	NotBefore   time.Time `bson:"NotBefore,omitempty" json:"NotBefore,omitempty"`
	NotAfter    time.Time `bson:"NotAfter,omitempty" json:"NotAfter,omitempty"`
}

type TimeStamp struct {
	*time.Time
}

type AppResponse struct {
	AppName               string    `json:"AppName"`
	CustomerName          string    `json:"CustomerName"`
	WebTechnology         int       `json:"WebTechnology"`
	Domain                string    `json:"Domain"`
	CallbackUrl           string    `json:"CallbackUrl"`
	DevDomain             string    `json:"DevDomain"`
	AppId                 int       `json:"AppId"`
	Key                   string    `json:"Key"`
	Secret                string    `json:"Secret"`
	Role                  string    `json:"Role"`
	IsWelcomeEmailEnabled bool      `json:"IsWelcomeEmailEnabled"`
	IsHttps               bool      `json:"IsHttps"`
	InterfaceId           int       `json:"InterfaceId"`
	RecurlyAccountCode    string    `json:"RecurlyAccountCode"`
	DateCreated           time.Time `json:"DateCreated"`
	DateModified          time.Time `json:"DateModified"`
	Status                bool      `json:"Status"`
	ProfilePhoto          string    `json:"ProfilePhoto"`
	ApiVersion            string    `json:"ApiVersion"`
	PrivacyPolicy         *string   `json:"PrivacyPolicy"`
	TermsOfService        *string   `json:"TermsOfService"`
	OwnerId               string    `json:"OwnerId"`
	ProductPlan           *string   `json:"ProductPlan"`
	EnvironmentTag        string    `json:"EnvironmentTag"`
}

type SendInvitation struct {
	Email   string   `json:"Email" validate:"required"`
	RoleIds []string `json:"RoleIds" validate:"required,min=1,dive,required"`
	// OrgId         string   `json:"-"`
	InvitationUrl string `json:"InvitationUrl"`
	InviterUid    string `json:"InviterUid" validate:"required"`
}

type UserRole struct {
	Id          string    `json:"Id"`
	RoleId      string    `json:"RoleId"`
	OrgId       string    `json:"OrgId"`
	Uid         string    `json:"Uid"`
	Name        string    `json:"Name"`
	Email       string    `json:"Email"`
	Username    string    `json:"Username"`
	CreatedDate time.Time `json:"CreatedDate"`
}

type RoleResponse struct {
	ID           string          `json:"Id"`
	Name         string          `json:"Name"`
	Description  string          `json:"Description"`
	Level        string          `json:"Level"`
	OrgId        string          `json:"OrgId"`
	Permissions  []PermissionB2B `json:"Permissions"`
	CreatedDate  time.Time       `bson:"CreatedDate" json:"CreatedDate"`
	ModifiedDate time.Time       `bson:"ModifiedDate" json:"ModifiedDate"`
}

type PermissionB2B struct {
	ID          string `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

type AllOrganizationsResponse struct {
	Id                       string      `json:"Id"`
	IsActive                 bool        `json:"IsActive"`
	Name                     string      `json:"Name"`
	Metadata                 interface{} `json:"Metadata"`
	IsAuthRestrictedToDomain bool        `json:"IsAuthRestrictedToDomain"`
	Policies                 struct {
		PasswordPolicy struct {
			MinLength          int  `json:"MinLength"`
			MaxLength          int  `json:"MaxLength"`
			RequireUppercase   bool `json:"RequireUppercase"`
			RequireLowercase   bool `json:"RequireLowercase"`
			RequireNumber      bool `json:"RequireNumber"`
			RequireSpecialChar bool `json:"RequireSpecialChar"`
			ExpiryDays         int  `json:"ExpiryDays"`
		} `json:"PasswordPolicy"`
		MFAPolicy struct {
			EnforcementMode string `json:"EnforcementMode"`
		} `json:"MFAPolicy"`
		SessionPolicy struct {
			AccessTokenTTL  int `json:"AccessTokenTTL"`
			RefreshTokenTTL int `json:"RefreshTokenTTL"`
		} `json:"SessionPolicy"`
	} `json:"Policies"`
	CreatedDate time.Time `json:"CreatedDate"`
}

type InvitationResponse struct {
	Data       []EmailVerificationResp `json:"Data"`
	TotalCount int64                   `json:"TotalCount"`
}

type EmailVerificationResp struct {
	Id      string    `json:"Id"`
	OrgId   string    `json:"OrgId,omitempty"`
	RoleIds []string  `json:"RoleIds,omitempty"`
	EMailId string    `json:"EmailId"`
	Date    time.Time `json:"CreatedDate"`
	Status  string    `json:"Status"`
}

// type EmailVerification struct {
// 	Id                string               `bson:"_id,omitempty" json:"Id"`
// 	IsVerified        bool                 `bson:"IsVerified" json:"-"`
// 	IsActive          bool                 `bson:"isActive" json:"-"`
// 	IsUsed            bool                 `bson:"IsUsed" json:"-"`
// 	AppId             *int                 `bson:"AppId,omitempty" json:"-"`
// 	Type              *uint8               `bson:"Type" json:"-"`
// 	EmailType         *string              `bson:"EmailType,omitempty" json:"-"`
// 	ClientGuid        *string              `bson:"ClientGuid,omitempty" json:"-"`
// 	Otp               *string              `bson:"Otp" json:"-"`
// 	GUID              *string              `bson:"GUID" json:"-"`
// 	Uid               *string              `bson:"Uid,omitempty" json:"-"`
// 	FK_UserId         string               `bson:"FK_UserId" json:"-"`
// 	EMailId           string               `bson:"EMailId" json:"EmailId"`
// 	Date              time.Time            `bson:"Date" json:"CreatedDate"`
// 	ExpirationDate    time.Time            `bson:"ExpirationDate" json:"ExpirationDate"`
// 	WorkflowSessionId *string              `bson:"wsid,omitempty" json:"-"`
// 	ClientId          *string              `bson:"cid,omitempty" json:"-"`
// 	IOSendEmailFlow   *bool                `bson:"iosendemailflow,omitempty" json:"-"`
// 	OrgId             primitive.ObjectID   `bson:"OrgId,omitempty" json:"OrgId"`
// 	RoleIds           []primitive.ObjectID `bson:"RoleIds,omitempty" json:"RoleIds"`
// 	ResentCount       int                  `bson:"ResentCount,omitempty" json:"-"` //Added if we need to track the number of times the email has been resent
// 	ModifiedDate      time.Time            `bson:"ModifiedDate,omitempty" json:"ModifiedDate"`
// 	InvitationUrl     string               `bson:"InvitationUrl,omitempty" json:"-"`
// 	InviterUid        string               `bson:"InviterUid,omitempty" json:"InviterUid"`
// 	InvitationToken   string               `bson:"InvitationToken,omitempty" json:"-"`
// }

type LoginResponse struct {
	organizationsList []OrganizationResponse `json:"organizationsList"`
	userId            string                 `json:"userId"`
	sessionToken      string                 `json:"sessionToken"`
}

type EmailValue struct {
	Type  string `json:"Type"`
	Value string `json:"Value"`
}
