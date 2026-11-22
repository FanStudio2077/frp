-- frp management console initial schema (SQLite)

CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	role TEXT NOT NULL DEFAULT 'viewer',
	disabled INTEGER NOT NULL DEFAULT 0,
	created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS servers (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	addr TEXT NOT NULL,
	port INTEGER NOT NULL,
	token_enc TEXT NOT NULL,
	tls_json TEXT,
	remark TEXT,
	created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS nodes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	token_hash TEXT NOT NULL,
	ip TEXT,
	os TEXT,
	arch TEXT,
	version TEXT,
	online INTEGER NOT NULL DEFAULT 0,
	last_seen TEXT,
	created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS instances (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	node_id INTEGER NOT NULL,
	server_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	admin_port INTEGER,
	status TEXT NOT NULL DEFAULT 'unknown',
	created_at TEXT NOT NULL DEFAULT (datetime('now')),
	UNIQUE(node_id, server_id)
);

CREATE TABLE IF NOT EXISTS proxies (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	instance_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	type TEXT NOT NULL,
	local_ip TEXT,
	local_port INTEGER,
	remote_port INTEGER,
	subdomain TEXT,
	custom_domains TEXT,
	extra_json TEXT,
	enabled INTEGER NOT NULL DEFAULT 1,
	created_at TEXT NOT NULL DEFAULT (datetime('now')),
	updated_at TEXT NOT NULL DEFAULT (datetime('now')),
	UNIQUE(instance_id, name)
);

CREATE TABLE IF NOT EXISTS visitors (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	instance_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	type TEXT NOT NULL,
	extra_json TEXT,
	enabled INTEGER NOT NULL DEFAULT 1,
	created_at TEXT NOT NULL DEFAULT (datetime('now')),
	UNIQUE(instance_id, name)
);

CREATE TABLE IF NOT EXISTS audit_logs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user TEXT,
	action TEXT NOT NULL,
	target TEXT,
	detail TEXT,
	ip TEXT,
	created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS alert_rules (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	metric TEXT NOT NULL,
	threshold INTEGER NOT NULL,
	scope TEXT NOT NULL DEFAULT 'global',
	enabled INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS alert_events (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	rule_id INTEGER,
	server_id INTEGER,
	value INTEGER,
	message TEXT,
	created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS settings (
	key TEXT PRIMARY KEY,
	value TEXT
);
