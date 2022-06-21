package utils

const (
	LocalBucket  = "pigil_storage"
	ConfigBucket = "pigil_config"
	CliDb        = "db"
	CliAuth      = "auth"
	CliChannels  = "channels"
	CliStatus    = "status"
	CliLogout    = "logout"
	CliHelp      = "help"
	CliDiscord   = "discord"
	DbFirstTime  = "isFirst"
	DbDiscordUrl = "discord_url"
	DbFalse      = "false"
	DbTrue       = "true"
	DatabaseName = "pigil.database"
	UserEmail    = "user_email"
	UserAT       = "access_token"
	UserRT       = "refresh_token"
	Discord      = "discord"
	Email        = "email"
)

var (
	GoogleClientId     = ""
	GoogleClientSecret = ""
	Channels           = []string{Email, Discord}
)
