package reality

import "fmt"

func (c Config) QuantumultX() (string, error) {
	uuid, err := c.UUID()
	if err != nil {
		return "", err
	}

	port, err := c.Port()
	if err != nil {
		return "", err
	}

	shortID, err := c.ShortID()
	if err != nil {
		return "", err
	}

	serverName, err := c.ServerName()
	if err != nil {
		return "", err
	}

	publicKey, err := c.PublicKey()
	if err != nil {
		return "", err
	}

	serverIP, err := publicIPv4()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"vless=%s:%d, method=none, password=%s, obfs=over-tls, obfs-host=%s, reality-base64-pubkey=%s, reality-hex-shortid=%s, vless-flow=xtls-rprx-vision, tag=Plachta-Reality",
		serverIP,
		port,
		uuid,
		serverName,
		publicKey,
		shortID,
	), nil
}
