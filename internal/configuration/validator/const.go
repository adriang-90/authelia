package validator

import (
	"regexp"

	"github.com/go-webauthn/webauthn/protocol"

	"github.com/authelia/authelia/v4/internal/oidc"
)

const (
	loopback           = "127.0.0.1"
	oauth2InstalledApp = "urn:ietf:wg:oauth:2.0:oob"
)

// Policy constants.
const (
	policyBypass    = "bypass"
	policyOneFactor = "one_factor"
	policyTwoFactor = "two_factor"
	policyDeny      = "deny"
)

// Hashing constants.
const (
	hashArgon2id = "argon2id"
	hashSHA512   = "sha512"
)

// Scheme constants.
const (
	schemeLDAP  = "ldap"
	schemeLDAPS = "ldaps"
	schemeHTTP  = "http"
	schemeHTTPS = "https"
)

// Test constants.
const (
	testInvalidPolicy = "invalid"
	testJWTSecret     = "a_secret"
	testLDAPBaseDN    = "base_dn"
	testLDAPPassword  = "password"
	testLDAPURL       = "ldap://ldap"
	testLDAPUser      = "user"
	testModeDisabled  = "disable"
	testEncryptionKey = "a_not_so_secure_encryption_key"
)

// Notifier Error constants.
const (
	errFmtNotifierMultipleConfigured = "notifier: please ensure only one of the 'smtp' or 'filesystem' notifier is configured"
	errFmtNotifierNotConfigured      = "notifier: you must ensure either the 'smtp' or 'filesystem' notifier " +
		"is configured"
	errFmtNotifierFileSystemFileNameNotConfigured = "notifier: filesystem: option 'filename' is required "
	errFmtNotifierSMTPNotConfigured               = "notifier: smtp: option '%s' is required"
)

// Authentication Backend Error constants.
const (
	errFmtAuthBackendNotConfigured = "authentication_backend: you must ensure either the 'file' or 'ldap' " +
		"authentication backend is configured"
	errFmtAuthBackendMultipleConfigured = "authentication_backend: please ensure only one of the 'file' or 'ldap' " +
		"backend is configured"
	errFmtAuthBackendRefreshInterval = "authentication_backend: option 'refresh_interval' is configured to '%s' but " +
		"it must be either a duration notation or one of 'disable', or 'always': %w"

	errFmtFileAuthBackendPathNotConfigured  = "authentication_backend: file: option 'path' is required"
	errFmtFileAuthBackendPasswordSaltLength = "authentication_backend: file: password: option 'salt_length' " +
		"must be 2 or more but it is configured a '%d'"
	errFmtFileAuthBackendPasswordUnknownAlg = "authentication_backend: file: password: option 'algorithm' " +
		"must be either 'argon2id' or 'sha512' but it is configured as '%s'"
	errFmtFileAuthBackendPasswordInvalidIterations = "authentication_backend: file: password: option " +
		"'iterations' must be 1 or more but it is configured as '%d'"
	errFmtFileAuthBackendPasswordArgon2idInvalidKeyLength = "authentication_backend: file: password: option " +
		"'key_length' must be 16 or more when using algorithm 'argon2id' but it is configured as '%d'"
	errFmtFileAuthBackendPasswordArgon2idInvalidParallelism = "authentication_backend: file: password: option " +
		"'parallelism' must be 1 or more when using algorithm 'argon2id' but it is configured as '%d'"
	errFmtFileAuthBackendPasswordArgon2idInvalidMemory = "authentication_backend: file: password: option 'memory' " +
		"must at least be parallelism multiplied by 8 when using algorithm 'argon2id' " +
		"with parallelism %d it should be at least %d but it is configured as '%d'"

	errFmtLDAPAuthBackendMissingOption = "authentication_backend: ldap: option '%s' is required"
	errFmtLDAPAuthBackendTLSMinVersion = "authentication_backend: ldap: tls: option " +
		"'minimum_tls_version' is invalid: %s: %w"
	errFmtLDAPAuthBackendImplementation = "authentication_backend: ldap: option 'implementation' " +
		"is configured as '%s' but must be one of the following values: '%s'"
	errFmtLDAPAuthBackendFilterReplacedPlaceholders = "authentication_backend: ldap: option " +
		"'%s' has an invalid placeholder: '%s' has been removed, please use '%s' instead"
	errFmtLDAPAuthBackendURLNotParsable = "authentication_backend: ldap: option " +
		"'url' could not be parsed: %w"
	errFmtLDAPAuthBackendURLInvalidScheme = "authentication_backend: ldap: option " +
		"'url' must have either the 'ldap' or 'ldaps' scheme but it is configured as '%s'"
	errFmtLDAPAuthBackendFilterEnclosingParenthesis = "authentication_backend: ldap: option " +
		"'%s' must contain enclosing parenthesis: '%s' should probably be '(%s)'"
	errFmtLDAPAuthBackendFilterMissingPlaceholder = "authentication_backend: ldap: option " +
		"'%s' must contain the placeholder '{%s}' but it is required"
)

// TOTP Error constants.
const (
	errFmtTOTPInvalidAlgorithm = "totp: option 'algorithm' must be one of '%s' but it is configured as '%s'"
	errFmtTOTPInvalidPeriod    = "totp: option 'period' option must be 15 or more but it is configured as '%d'"
	errFmtTOTPInvalidDigits    = "totp: option 'digits' must be 6 or 8 but it is configured as '%d'"
)

// Storage Error constants.
const (
	errStrStorage                            = "storage: configuration for a 'local', 'mysql' or 'postgres' database must be provided"
	errStrStorageEncryptionKeyMustBeProvided = "storage: option 'encryption_key' must is required"
	errStrStorageEncryptionKeyTooShort       = "storage: option 'encryption_key' must be 20 characters or longer"
	errFmtStorageUserPassMustBeProvided      = "storage: %s: option 'username' and 'password' are required" //nolint: gosec
	errFmtStorageOptionMustBeProvided        = "storage: %s: option '%s' is required"
	errFmtStoragePostgreSQLInvalidSSLMode    = "storage: postgres: ssl: option 'mode' must be one of '%s' but it is configured as '%s'"
)

// OpenID Error constants.
const (
	errFmtOIDCNoClientsConfigured = "identity_providers: oidc: option 'clients' must have one or " +
		"more clients configured"
	errFmtOIDCNoPrivateKey = "identity_providers: oidc: option 'issuer_private_key' is required"

	errFmtOIDCEnforcePKCEInvalidValue = "identity_providers: oidc: option 'enforce_pkce' must be 'never', " +
		"'public_clients_only' or 'always', but it is configured as '%s'"

	errFmtOIDCClientsDuplicateID = "identity_providers: oidc: one or more clients have the same id but all client" +
		"id's must be unique"
	errFmtOIDCClientsWithEmptyID = "identity_providers: oidc: one or more clients have been configured with " +
		"an empty id"

	errFmtOIDCClientInvalidSecret       = "identity_providers: oidc: client '%s': option 'secret' is required"
	errFmtOIDCClientPublicInvalidSecret = "identity_providers: oidc: client '%s': option 'secret' is " +
		"required to be empty when option 'public' is true"
	errFmtOIDCClientRedirectURI = "identity_providers: oidc: client '%s': option 'redirect_uris' has an " +
		"invalid value: redirect uri '%s' must have a scheme of 'http' or 'https' but '%s' is configured"
	errFmtOIDCClientRedirectURICantBeParsed = "identity_providers: oidc: client '%s': option 'redirect_uris' has an " +
		"invalid value: redirect uri '%s' could not be parsed: %v"
	errFmtOIDCClientRedirectURIPublic = "identity_providers: oidc: client '%s': option 'redirect_uris' has the" +
		"redirect uri '%s' when option 'public' is false but this is invalid as this uri is not valid " +
		"for the openid connect confidential client type"
	errFmtOIDCClientRedirectURIAbsolute = "identity_providers: oidc: client '%s': option 'redirect_uris' has an " +
		"invalid value: redirect uri '%s' must have the scheme 'http' or 'https' but it has no scheme"
	errFmtOIDCClientInvalidPolicy = "identity_providers: oidc: client '%s': option 'policy' must be 'one_factor' " +
		"or 'two_factor' but it is configured as '%s'"
	errFmtOIDCClientInvalidEntry = "identity_providers: oidc: client '%s': option '%s' must only have the values " +
		"'%s' but one option is configured as '%s'"
	errFmtOIDCClientInvalidUserinfoAlgorithm = "identity_providers: oidc: client '%s': option " +
		"'userinfo_signing_algorithm' must be one of '%s' but it is configured as '%s'"
	errFmtOIDCServerInsecureParameterEntropy = "openid connect provider: SECURITY ISSUE - minimum parameter entropy is " +
		"configured to an unsafe value, it should be above 8 but it's configured to %d"
)

// Webauthn Error constants.
const (
	errFmtWebauthnConveyancePreference = "webauthn: option 'attestation_conveyance_preference' must be one of '%s' but it is configured as '%s'"
	errFmtWebauthnUserVerification     = "webauthn: option 'user_verification' must be one of 'discouraged', 'preferred', 'required' but it is configured as '%s'"
)

// Access Control error constants.
const (
	errFmtAccessControlDefaultPolicyValue = "access control: option 'default_policy' must be one of '%s' but it is " +
		"configured as '%s'"
	errFmtAccessControlDefaultPolicyWithoutRules = "access control: 'default_policy' option '%s' is invalid: when " +
		"no rules are specified it must be 'two_factor' or 'one_factor'"
	errFmtAccessControlNetworkGroupIPCIDRInvalid = "access control: networks: network group '%s' is invalid: the " +
		"network '%s' is not a valid IP or CIDR notation"
	errFmtAccessControlWarnNoRulesDefaultPolicy = "access control: no rules have been specified so the " +
		"'default_policy' of '%s' is going to be applied to all requests"
	errFmtAccessControlRuleNoDomains = "access control: rule %s: rule is invalid: must have the option " +
		"'domain' configured"
	errFmtAccessControlRuleInvalidPolicy = "access control: rule %s: rule 'policy' option '%s' " +
		"is invalid: must be one of 'deny', 'two_factor', 'one_factor' or 'bypass'"
	errAccessControlRuleBypassPolicyInvalidWithSubjects = "access control: rule %s: 'policy' option 'bypass' is " +
		"not supported when 'subject' option is configured: see " +
		"https://www.authelia.com/docs/configuration/access-control.html#bypass"
	errFmtAccessControlRuleNetworksInvalid = "access control: rule %s: the network '%s' is not a " +
		"valid Group Name, IP, or CIDR notation"
	errFmtAccessControlRuleResourceInvalid = "access control: rule %s: 'resources' option '%s' is " +
		"invalid: %w"
	errFmtAccessControlRuleSubjectInvalid = "access control: rule %s: 'subject' option '%s' is " +
		"invalid: must start with 'user:' or 'group:'"
	errFmtAccessControlRuleMethodInvalid = "access control: rule %s: 'methods' option '%s' is " +
		"invalid: must be one of '%s'"
)

// Theme Error constants.
const (
	errFmtThemeName = "option 'theme' must be one of '%s' but it is configured as '%s'"
)

// NTP Error constants.
const (
	errFmtNTPVersion = "ntp: option 'version' must be either 3 or 4 but it is configured as '%d'"
)

// Session error constants.
const (
	errFmtSessionOptionRequired           = "session: option '%s' is required"
	errFmtSessionDomainMustBeRoot         = "session: option 'domain' must be the domain you wish to protect not a wildcard domain but it is configured as '%s'"
	errFmtSessionSameSite                 = "session: option 'same_site' must be one of '%s' but is configured as '%s'"
	errFmtSessionSecretRequired           = "session: option 'secret' is required when using the '%s' provider"
	errFmtSessionRedisPortRange           = "session: redis: option 'port' must be between 1 and 65535 but is configured as '%d'"
	errFmtSessionRedisHostRequired        = "session: redis: option 'host' is required"
	errFmtSessionRedisHostOrNodesRequired = "session: redis: option 'host' or the 'high_availability' option 'nodes' is required"

	errFmtSessionRedisSentinelMissingName     = "session: redis: high_availability: option 'sentinel_name' is required"
	errFmtSessionRedisSentinelNodeHostMissing = "session: redis: high_availability: option 'nodes': option 'host' is required for each node but one or more nodes are missing this"
)

// Regulation Error Consts.
const (
	errFmtRegulationFindTimeGreaterThanBanTime = "regulation: option 'find_time' must be less than or equal to option 'ban_time'"
)

// Server Error constants.
const (
	errFmtServerTLSCert                           = "server: tls: option 'key' must also be accompanied by option 'certificate'"
	errFmtServerTLSKey                            = "server: tls: option 'certificate' must also be accompanied by option 'key'"
	errFmtServerTLSCertFileDoesNotExist           = "server: tls: file path %s provided in 'certificate' does not exist"
	errFmtServerTLSKeyFileDoesNotExist            = "server: tls: file path %s provided in 'key' does not exist"
	errFmtServerTLSClientAuthCertFileDoesNotExist = "server: tls: client_certificates: certificates: file path %s does not exist"
	errFmtServerTLSClientAuthNoAuth               = "server: tls: client authentication cannot be configured if no server certificate and key are provided"

	errFmtServerPathNoForwardSlashes = "server: option 'path' must not contain any forward slashes"
	errFmtServerPathAlphaNum         = "server: option 'path' must only contain alpha numeric characters"
	errFmtServerBufferSize           = "server: option '%s_buffer_size' must be above 0 but it is configured as '%d'"
)

// Error constants.
const (
	/*
		errFmtDeprecatedConfigurationKey = "the %s configuration option is deprecated and will be " +
			"removed in %s, please use %s instead"

		Uncomment for use when deprecating keys.

		TODO: Create a method from within Koanf to automatically remap deprecated keys and produce warnings.
		TODO (cont): The main consideration is making sure we do not overwrite the destination key name if it already exists.
	*/

	errFmtReplacedConfigurationKey = "invalid configuration key '%s' was replaced by '%s'"

	errFmtLoggingLevelInvalid = "log: option 'level' must be one of '%s' but it is configured as '%s'"

	errFileHashing  = "config key incorrect: authentication_backend.file.hashing should be authentication_backend.file.password"
	errFilePHashing = "config key incorrect: authentication_backend.file.password_hashing should be authentication_backend.file.password"
	errFilePOptions = "config key incorrect: authentication_backend.file.password_options should be authentication_backend.file.password"
)

var validStoragePostgreSQLSSLModes = []string{testModeDisabled, "require", "verify-ca", "verify-full"}

var validThemeNames = []string{"light", "dark", "grey", "auto"}

var validSessionSameSiteValues = []string{"none", "lax", "strict"}

var validLoLevels = []string{"trace", "debug", "info", "warn", "error"}

var validWebauthnConveyancePreferences = []string{string(protocol.PreferNoAttestation), string(protocol.PreferIndirectAttestation), string(protocol.PreferDirectAttestation)}
var validWebauthnUserVerificationRequirement = []string{string(protocol.VerificationDiscouraged), string(protocol.VerificationPreferred), string(protocol.VerificationRequired)}

var validACLRuleMethods = []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "TRACE", "CONNECT", "OPTIONS"}
var validACLRulePolicies = []string{policyBypass, policyOneFactor, policyTwoFactor, policyDeny}

var validOIDCScopes = []string{oidc.ScopeOpenID, oidc.ScopeEmail, oidc.ScopeProfile, oidc.ScopeGroups, "offline_access"}
var validOIDCGrantTypes = []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"}
var validOIDCResponseModes = []string{"form_post", "query", "fragment"}
var validOIDCUserinfoAlgorithms = []string{"none", "RS256"}

var reKeyReplacer = regexp.MustCompile(`\[\d+]`)

// ValidKeys is a list of valid keys that are not secret names. For the sake of consistency please place any secret in
// the secret names map and reuse it in relevant sections.
var ValidKeys = []string{
	// Root Keys.
	"certificates_directory",
	"theme",
	"default_redirection_url",
	"jwt_secret",

	// Log keys.
	"log.level",
	"log.format",
	"log.file_path",
	"log.keep_stdout",

	// Server Keys.
	"server.host",
	"server.port",
	"server.read_buffer_size",
	"server.write_buffer_size",
	"server.path",
	"server.asset_path",
	"server.enable_pprof",
	"server.enable_expvars",
	"server.disable_healthcheck",
	"server.tls.key",
	"server.tls.certificate",
	"server.headers.csp_template",

	// TOTP Keys.
	"totp.disable",
	"totp.issuer",
	"totp.algorithm",
	"totp.digits",
	"totp.period",
	"totp.skew",

	// Webauthn Keys.
	"webauthn.disable",
	"webauthn.display_name",
	"webauthn.attestation_conveyance_preference",
	"webauthn.user_verification",
	"webauthn.timeout",

	// DUO API Keys.
	"duo_api.hostname",
	"duo_api.enable_self_enrollment",
	"duo_api.secret_key",
	"duo_api.integration_key",

	// Access Control Keys.
	"access_control.default_policy",
	"access_control.networks",
	"access_control.networks[].name",
	"access_control.networks[].networks",
	"access_control.rules",
	"access_control.rules[].domain",
	"access_control.rules[].methods",
	"access_control.rules[].networks",
	"access_control.rules[].subject",
	"access_control.rules[].policy",
	"access_control.rules[].resources",

	// Session Keys.
	"session.name",
	"session.domain",
	"session.secret",
	"session.same_site",
	"session.expiration",
	"session.inactivity",
	"session.remember_me_duration",

	// Redis Session Keys.
	"session.redis.host",
	"session.redis.port",
	"session.redis.username",
	"session.redis.password",
	"session.redis.database_index",
	"session.redis.maximum_active_connections",
	"session.redis.minimum_idle_connections",
	"session.redis.tls.minimum_version",
	"session.redis.tls.skip_verify",
	"session.redis.tls.server_name",
	"session.redis.high_availability.sentinel_name",
	"session.redis.high_availability.sentinel_username",
	"session.redis.high_availability.sentinel_password",
	"session.redis.high_availability.nodes",
	"session.redis.high_availability.nodes[].host",
	"session.redis.high_availability.nodes[].port",
	"session.redis.high_availability.route_by_latency",
	"session.redis.high_availability.route_randomly",

	// Storage Keys.
	"storage.encryption_key",

	// Local Storage Keys.
	"storage.local.path",

	// MySQL Storage Keys.
	"storage.mysql.host",
	"storage.mysql.port",
	"storage.mysql.database",
	"storage.mysql.username",
	"storage.mysql.password",
	"storage.mysql.timeout",

	// PostgreSQL Storage Keys.
	"storage.postgres.host",
	"storage.postgres.port",
	"storage.postgres.database",
	"storage.postgres.username",
	"storage.postgres.password",
	"storage.postgres.timeout",
	"storage.postgres.schema",
	"storage.postgres.ssl.mode",
	"storage.postgres.ssl.root_certificate",
	"storage.postgres.ssl.certificate",
	"storage.postgres.ssl.key",

	"storage.postgres.sslmode", // Deprecated. TODO: Remove in v4.36.0.

	// FileSystem Notifier Keys.
	"notifier.filesystem.filename",
	"notifier.disable_startup_check",

	// SMTP Notifier Keys.
	"notifier.smtp.host",
	"notifier.smtp.port",
	"notifier.smtp.timeout",
	"notifier.smtp.username",
	"notifier.smtp.password",
	"notifier.smtp.identifier",
	"notifier.smtp.sender",
	"notifier.smtp.subject",
	"notifier.smtp.startup_check_address",
	"notifier.smtp.disable_require_tls",
	"notifier.smtp.disable_html_emails",
	"notifier.smtp.tls.minimum_version",
	"notifier.smtp.tls.skip_verify",
	"notifier.smtp.tls.server_name",

	// Regulation Keys.
	"regulation.max_retries",
	"regulation.find_time",
	"regulation.ban_time",

	// Authentication Backend Keys.
	"authentication_backend.disable_reset_password",
	"authentication_backend.refresh_interval",

	// LDAP Authentication Backend Keys.
	"authentication_backend.ldap.implementation",
	"authentication_backend.ldap.url",
	"authentication_backend.ldap.timeout",
	"authentication_backend.ldap.base_dn",
	"authentication_backend.ldap.username_attribute",
	"authentication_backend.ldap.additional_users_dn",
	"authentication_backend.ldap.users_filter",
	"authentication_backend.ldap.additional_groups_dn",
	"authentication_backend.ldap.groups_filter",
	"authentication_backend.ldap.group_name_attribute",
	"authentication_backend.ldap.mail_attribute",
	"authentication_backend.ldap.display_name_attribute",
	"authentication_backend.ldap.user",
	"authentication_backend.ldap.password",
	"authentication_backend.ldap.start_tls",
	"authentication_backend.ldap.tls.minimum_version",
	"authentication_backend.ldap.tls.skip_verify",
	"authentication_backend.ldap.tls.server_name",

	// File Authentication Backend Keys.
	"authentication_backend.file.path",
	"authentication_backend.file.password.algorithm",
	"authentication_backend.file.password.iterations",
	"authentication_backend.file.password.key_length",
	"authentication_backend.file.password.salt_length",
	"authentication_backend.file.password.memory",
	"authentication_backend.file.password.parallelism",

	// Identity Provider Keys.
	"identity_providers.oidc.hmac_secret",
	"identity_providers.oidc.issuer_private_key",
	"identity_providers.oidc.id_token_lifespan",
	"identity_providers.oidc.access_token_lifespan",
	"identity_providers.oidc.refresh_token_lifespan",
	"identity_providers.oidc.authorize_code_lifespan",
	"identity_providers.oidc.enforce_pkce",
	"identity_providers.oidc.enable_pkce_plain_challenge",
	"identity_providers.oidc.enable_client_debug_messages",
	"identity_providers.oidc.minimum_parameter_entropy",
	"identity_providers.oidc.clients",
	"identity_providers.oidc.clients[].id",
	"identity_providers.oidc.clients[].description",
	"identity_providers.oidc.clients[].public",
	"identity_providers.oidc.clients[].secret",
	"identity_providers.oidc.clients[].redirect_uris",
	"identity_providers.oidc.clients[].authorization_policy",
	"identity_providers.oidc.clients[].scopes",
	"identity_providers.oidc.clients[].audience",
	"identity_providers.oidc.clients[].grant_types",
	"identity_providers.oidc.clients[].response_types",
	"identity_providers.oidc.clients[].response_modes",
	"identity_providers.oidc.clients[].userinfo_signing_algorithm",

	// NTP keys.
	"ntp.address",
	"ntp.version",
	"ntp.max_desync",
	"ntp.disable_startup_check",
	"ntp.disable_failure",
}

var replacedKeys = map[string]string{
	"authentication_backend.ldap.skip_verify":         "authentication_backend.ldap.tls.skip_verify",
	"authentication_backend.ldap.minimum_tls_version": "authentication_backend.ldap.tls.minimum_version",
	"notifier.smtp.disable_verify_cert":               "notifier.smtp.tls.skip_verify",
	"logs_level":                                      "log.level",
	"logs_file_path":                                  "log.file_path",
	"log_level":                                       "log.level",
	"log_file_path":                                   "log.file_path",
	"log_format":                                      "log.format",
	"host":                                            "server.host",
	"port":                                            "server.port",
	"tls_key":                                         "server.tls.key",
	"tls_cert":                                        "server.tls.certificate",
}

var specificErrorKeys = map[string]string{
	"google_analytics": "config key removed: google_analytics - this functionality has been deprecated",
	"notifier.smtp.trusted_cert": "invalid configuration key 'notifier.smtp.trusted_cert' it has been removed, " +
		"option has been replaced by the global option 'certificates_directory'",

	"authentication_backend.file.password_options.algorithm":   errFilePOptions,
	"authentication_backend.file.password_options.iterations":  errFilePOptions,
	"authentication_backend.file.password_options.key_length":  errFilePOptions,
	"authentication_backend.file.password_options.salt_length": errFilePOptions,
	"authentication_backend.file.password_options.memory":      errFilePOptions,
	"authentication_backend.file.password_options.parallelism": errFilePOptions,
	"authentication_backend.file.password_hashing.algorithm":   errFilePHashing,
	"authentication_backend.file.password_hashing.iterations":  errFilePHashing,
	"authentication_backend.file.password_hashing.key_length":  errFilePHashing,
	"authentication_backend.file.password_hashing.salt_length": errFilePHashing,
	"authentication_backend.file.password_hashing.memory":      errFilePHashing,
	"authentication_backend.file.password_hashing.parallelism": errFilePHashing,
	"authentication_backend.file.hashing.algorithm":            errFileHashing,
	"authentication_backend.file.hashing.iterations":           errFileHashing,
	"authentication_backend.file.hashing.key_length":           errFileHashing,
	"authentication_backend.file.hashing.salt_length":          errFileHashing,
	"authentication_backend.file.hashing.memory":               errFileHashing,
	"authentication_backend.file.hashing.parallelism":          errFileHashing,
}
