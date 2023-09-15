package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/OpenPajak/efaktur-go/pkg/provider/web"
	"github.com/ii64/go-binder/binder"
	"github.com/ii64/go-binder/binder/ext/json"
	"github.com/ii64/go-binder/binder/ext/toml"
	"github.com/ii64/go-binder/binder/ext/yaml"
	_ "github.com/joho/godotenv/autoload"
)

var configPath = "efaktur.toml"

type Config struct {
	CertificatePath     string `bind:"certPath" json:"certPath" yaml:"certPath" toml:"certPath" argx:"certPath" env:"CERT_PATH" usage:"PKCS#12 certificate path"`
	CertificatePassword string `bind:"certPasswd" json:"certPasswd"  yaml:"certPasswd" toml:"certPasswd" argx:"certPasswd" env:"CERT_PASS" usage:"PKCS#12 certificate password"`
	ServicePassword     string `bind:"svcPasswd" json:"svcPasswd" yaml:"svcPasswd" toml:"svcPasswd" argx:"svcPasswd" env:"SVC_PASSWD" usage:"Service password"`

	DebugMode bool `json:"-" yaml:"-" toml:"-" argx:"debug" env:"DEBUG" usage:"Debug mode"`
}

var conf = &Config{
	CertificatePath:     "",
	CertificatePassword: "",
	ServicePassword:     "",
	DebugMode:           false,
}

func init() {
	if envConfigPath := os.Getenv("CONFIG"); envConfigPath != "" {
		configPath = envConfigPath
	}

	switch filepath.Ext(configPath) {
	case ".json":
		binder.LoadConfig = json.LoadConfig(configPath)
		binder.SaveConfig = json.SaveConfig(configPath, "  ")
	case ".yaml", ".yml":
		binder.LoadConfig = yaml.LoadConfig(configPath)
		binder.SaveConfig = yaml.SaveConfig(configPath, 2)
	case ".toml":
		binder.LoadConfig = toml.LoadConfig(configPath)
		binder.SaveConfig = toml.SaveConfig(configPath, "  ")
	}
	binder.SaveOnClose = true

	// bind
	binder.BindArgsConf(conf, "conf")
}

var cmds = Commands{
	cmdDump,
	cmdPrepopulated,
}

func main() {
retryGobinderInit:
	retryAttempt := 0
	if retryAttempt > 10 {
		panic(fmt.Sprintf("go-binder init failed"))
	}
	if err := binder.Init(); err != nil {
		// Can't use [`os.IsNotExist`] because it does not support Unwrap.
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, io.EOF) {
			binder.In()
			if err = binder.Save(); err != nil {
				panic(err)
			}
			// retry init after saving config file.
			retryAttempt++
			goto retryGobinderInit
		} else {
			panic(err)
		}
	}
	flag.Parse()

	binder.In()
	defer binder.Close()

	log.SetOutput(os.Stderr)
	if !conf.DebugMode {
		log.SetOutput(io.Discard)
	}

	args := flag.Args()
	if len(args) < 1 || cmds.Lookup(args[0]) == nil {
		cmds.Usage()
	}

	Main(args)
}

func Main(args []string) {
	tlsCert, clientCAs, err := web.PKCS12ToTLSCertificateFromFile(conf.CertificatePath, conf.CertificatePassword)
	if err != nil {
		log.Fatalf(
			"failed to open PKCS#12 certificate from file: %q: %s\n",
			conf.CertificatePath,
			err)
	}

	log.Printf("[CERT-LEAF] -----------------------------\n")
	log.Printf("[CERT-LEAF] SN  : %s\n", tlsCert.Leaf.SerialNumber)
	log.Printf("[CERT-LEAF] ISS : %q\n", tlsCert.Leaf.Issuer)
	log.Printf("[CERT-LEAF] SUB : %q\n", tlsCert.Leaf.Subject)

	expiryDur := tlsCert.Leaf.NotAfter.Sub(time.Now().UTC())
	log.Printf("[CERT-LEAF] NOT BEFORE: %q", tlsCert.Leaf.NotBefore)
	log.Printf("[CERT-LEAF] EXP IN %s [%.2f day(s)] (on %q)",
		expiryDur,
		expiryDur.Hours()/24,
		tlsCert.Leaf.NotAfter)
	log.Printf("[CERT-LEAF] -----------------------------\n")

	if expiryDur < (time.Hour * 24 * 7 * 2) {
		log.Printf("[WARNING] LEAF CERTIFICATE IS ABOUT TO EXPIRED IN %.2f DAY(s) !!!\n",
			expiryDur.Hours()/24)
	} else if expiryDur < 0 {
		log.Fatalf("Leaf certificate is expired\n")
	}

	client, err := web.NewClient(web.ClientOptions{
		TLSCertificate: tlsCert,
		TLSClientCAs:   clientCAs,

		// Temporarily
		TLSInsecureSkipVerify: true,
	})
	if err != nil {
		log.Fatalf("failed to create client: %s\n", err)
	}

	ctx := context.Background()
	ctx = WithConf(ctx, conf)
	ctx = WithClient(ctx, client)

	switch err := cmds.Run(ctx, args); err {
	case flag.ErrHelp:
		os.Exit(1)
	case ErrCmdNotExist:
		goto cmdUnavail
	default:
		if err != nil {
			panic(err)
		}
	}
	return

cmdUnavail:
	cmds.Usage()
}

func cmdPrerun(ctx context.Context) func() {
	conf := GetConfFromContext(ctx)
	client := GetClientFromContext(ctx)

	if err := client.Login(ctx, conf.ServicePassword); err != nil {
		log.Fatalf("failed to login: %s\n", err)
	}
	log.Printf("logged in")

	profileResponse, err := client.Profile.Get(ctx)
	if err != nil {
		log.Fatalf("failed to get PKP profile: %s\n", err)
	}
	profile := profileResponse.GetOne()
	if profile == nil {
		log.Fatalf("no PKP profile entry\n")
	}
	log.Printf(
		"profile npwp=%q name=%q",
		profile.Npwp,
		profile.Nama,
	)

	return func() {
		if err := client.Logout(ctx); err != nil {
			log.Fatalf("failed to logout: %s\n", err)
		}
		log.Printf("logged out")
	}
}

type Commands []Cmd

var ErrCmdNotExist = errors.New("cmd not exist")

func (c Commands) Usage() {
	c.PrintAvailable()
	os.Exit(1)
}

func (c Commands) PrintAvailable() {
	var sb strings.Builder
	sb.WriteString("available commands:\n")
	for i, cmd := range c {
		sb.WriteString(" - ")
		sb.WriteString(cmd.Name())
		if i+1 < len(c) {
			sb.WriteRune('\n')
		}
	}
	fmt.Println(sb.String())
}

func (c Commands) Lookup(name string) Cmd {
	for _, cmd := range c {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}

func (c Commands) Run(ctx context.Context, args []string) (err error) {
	if len(args) < 1 {
		return nil
	}
	cmd := c.Lookup(args[0])
	if cmd == nil {
		return ErrCmdNotExist
	}
	err = cmd.Run(ctx, args[1:])
	return
}

type Cmd interface {
	Name() string
	Run(ctx context.Context, args []string) (err error)
}

type cmd[T any] struct {
	name     string
	fs       *flag.FlagSet
	data     T
	setup    func(fs *flag.FlagSet, data *T)
	callback func(ctx context.Context, c *cmd[T]) (err error)
}

func (c *cmd[T]) init() *cmd[T] {
	c.fs = flag.NewFlagSet(c.name, flag.ContinueOnError)
	if f := c.setup; f != nil {
		c.setup(c.fs, &c.data)
	}
	return c
}

func (c *cmd[T]) Name() string {
	return c.name
}

func (c *cmd[T]) Run(ctx context.Context, args []string) (err error) {
	if c.fs == nil {
		c.init()
	}
	if err = c.fs.Parse(args); err != nil {
		return
	}
	if f := c.callback; f != nil {
		return f(ctx, c)
	}
	return
}
