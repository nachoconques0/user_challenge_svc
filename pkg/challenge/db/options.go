package db

// Retrieve the default options
func defaultOptions() Options {
	return Options{
		Host:          "127.0.0.1",
		Port:          "5432",
		User:          "postgres",
		Password:      "postgres",
		Database:      "postgres",
		SingularTable: true,
		SSLMode:       "disable",
	}
}

type Options struct {
	Host           string
	Port           string
	User           string
	Password       string
	Database       string
	MaxConnections int
	SingularTable  bool
	SSLMode        string
	Debug          bool
}

// WithHost takes the host of the database we want to connect to
func WithHost(h string) Option {
	return func(o *Options) {
		o.Host = h
	}
}

// WithPort takes the port of the database we want to connect to
func WithPort(p string) Option {
	return func(o *Options) {
		o.Port = p
	}
}

// WithDatabase takes the database name we want to connect to
func WithDatabase(d string) Option {
	return func(o *Options) {
		o.Database = d
	}
}

// WithUser takes the user of the database we want to connect to
func WithUser(u string) Option {
	return func(o *Options) {
		o.User = u
	}
}

// WithPassword takes the user's password of the database we want to connect to
func WithPassword(p string) Option {
	return func(o *Options) {
		o.Password = p
	}
}

// WithMaxConnections sets a limit on the maximum number of open connections
func WithMaxConnections(maxC int) Option {
	return func(o *Options) {
		o.MaxConnections = maxC
	}
}

// WithSingularTable sets if the table names should be singular
func WithSingularTable(s bool) Option {
	return func(o *Options) {
		o.SingularTable = s
	}
}

// WithSSLMode sets the SSL mode for the connection
func WithSSLMode(mode string) Option {
	return func(o *Options) {
		o.SSLMode = mode
	}
}

// WithDebug controls if queries should be logged
func WithDebug(d bool) Option {
	return func(o *Options) {
		o.Debug = d
	}
}

type Option func(*Options)
