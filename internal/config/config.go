package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken string
	ChatGptApiKey string

	Repo      Repo
	Promo     Promo
	Payment   Payment
	Tesseract Tesseract
	ChatAI    ChatAI
	Bot       Bot
	Server    Server
}

type Repo struct {
	UsersBasePath string `mapstructure:"users_base_path"`
	UrlsBasePath  string `mapstructure:"urls_base_path"`
}

type Promo struct {
	BotName        string `mapstructure:"bot_name"`
	Gift           string `mapstructure:"gift"`
	Start          string `mapstructure:"start"`
	CountReqFriend int64  `mapstructure:"count_req_friend"`
	CountReqAuthor int64  `mapstructure:"count_req_author"`
}

type Payment struct {
	Wallet     string
	KeyPayment string // 32 byte string
	Message    string `mapstructure:"message"`
	URL        URL
	Products   Shop
}

type URL struct {
	Scheme string `mapstructure:"scheme"`
	Host   string `mapstructure:"host"`
	Path   string `mapstructure:"path"`
}

type Tesseract struct {
	NlpLanguages []string `mapstructure:"nlp_languages"`
}

type ChatAI struct {
	RoleContent string `mapstructure:"role_content"`
}

type Bot struct {
	Ip   string `mapstructure:"ip"`
	Port uint   `mapstructure:"port"`
	Url  string `mapstructure:"url"`

	CertPath string `mapstructure:"cert_path"`
	KeyPath  string `mapstructure:"key_path"`

	StartBalance int64   `mapstructure:"start_balance"`
	AdminsId     []int64 `mapstructure:"admins_id"`
	Messages     Messages
	Keyboard     Keyboard
	Commands     Commands
}

type Shop struct {
	ProductPrice01      int64 `mapstructure:"product_price_01"`
	ProductPrice02      int64 `mapstructure:"product_price_02"`
	ProductPrice03      int64 `mapstructure:"product_price_03"`
	ProductSalesPrice01 int64 `mapstructure:"product_sales_price_01"`
	ProductSalesPrice02 int64 `mapstructure:"product_sales_price_02"`
	ProductSalesPrice03 int64 `mapstructure:"product_sales_price_03"`
	ProductBasePrice01  int64 `mapstructure:"product_base_price_01"`
	ProductBasePrice02  int64 `mapstructure:"product_base_price_02"`
	ProductBasePrice03  int64 `mapstructure:"product_base_price_03"`
	ProductCount01      int64 `mapstructure:"product_count_01"`
	ProductCount02      int64 `mapstructure:"product_count_02"`
	ProductCount03      int64 `mapstructure:"product_count_03"`
	SalesCount          int64 `mapstructure:"sales_count"`
}

type Messages struct {
	Responses
	Errors
}

type Responses struct {
	Start            string `mapstructure:"start"`
	Help             string `mapstructure:"help"`
	BuyRequests      string `mapstructure:"buy_requests"`
	BuyRequestsSales string `mapstructure:"buy_requests_sales"`
	Balance          string `mapstructure:"balance"`
	Subscribe        string `mapstructure:"subscribe"`
	Image            string `mapstructure:"image"`
	Task             string `mapstructure:"task"`
	TextGot          string `mapstructure:"got_text"`
	ImageGot         string `mapstructure:"got_image"`
	Admin            string `mapstructure:"admin"`
	AdminDbUsers     string `mapstructure:"admin_db_users"`
	AdminDbUrls      string `mapstructure:"admin_db_urls"`
	AdviceButton     string `mapstructure:"advice_button"`
	AdviceStart      string `mapstructure:"advice_start"`
	Advice01         string `mapstructure:"advice01"`
	Advice02         string `mapstructure:"advice02"`
	Advice03         string `mapstructure:"advice03"`
	Advice04         string `mapstructure:"advice04"`
	Promo            string `mapstructure:"promo"`
	PromoOK          string `mapstructure:"promo_ok"`
	CreateRef01      string `mapstructure:"create_ref01"`
	CreateRef02      string `mapstructure:"create_ref02"`
	CreateRef03      string `mapstructure:"create_ref03"`
}

type Errors struct {
	UnknownMessage  string `mapstructure:"unknown_message"`
	UnknownCommand  string `mapstructure:"unknown_command"`
	UnknownCallback string `mapstructure:"unknown_callback"`
	NotAuth         string `mapstructure:"not_auth"`
	AdminError      string `mapstructure:"admin_error"`
	AdminUsageError string `mapstructure:"admin_usage_error"`
	PromoError      string `mapstructure:"promo_error"`
	PromoUsedError  string `mapstructure:"promo_used_error"`
	Default         string `mapstructure:"default"`
}

type Keyboard struct {
	Menu    Menu
	Balance Balance
	Advices Advices
}

type Menu struct {
	Help           string `mapstructure:"help"`
	Balance        string `mapstructure:"balance"`
	Ref            string `mapstructure:"ref"`
	ImageRecognize string `mapstructure:"image_recognize"`
	SolveTask      string `mapstructure:"solve_task"`
}

type Balance struct {
	Buy string `mapstructure:"buy"`
	Ref string `mapstructure:"ref"`
}

type Advices struct {
	Advice string `mapstructure:"Advice"`
}

type Commands struct {
	Start        string `mapstructure:"start"`
	Help         string `mapstructure:"help"`
	Admin        string `mapstructure:"admin"`
	Logs         string `mapstructure:"logs"`
	DbUsers      string `mapstructure:"db_users"`
	DbUrls       string `mapstructure:"db_urls"`
	CreateUser   string `mapstructure:"create_user"`
	CreateUrl    string `mapstructure:"create_url"`
	DeleteUser   string `mapstructure:"delete_user"`
	DeleteUrl    string `mapstructure:"delete_url"`
	SendAll      string `mapstructure:"send_all"`
	SendAllBuy   string `mapstructure:"send_all_buy"`
	SendAllRef   string `mapstructure:"send_all_ref"`
	SendAllZeros string `mapstructure:"send_all_zeros"`
	SetSales     string `mapstructure:"set_sales"`
	GetSales     string `mapstructure:"get_sales"`
	GetTokens    string `mapstructure:"get_tokens"`
	AddToken     string `mapstructure:"add_tokens"`
	RemoveToken  string `mapstructure:"remove_tokens"`
	NextToken    string `mapstructure:"next_tokens"`
	SendAllAd    string `mapstructure:"send_all_ad"`
}

type Server struct {
	Port uint `mapstructure:"port"`
}

func Init() (*Config, error) {
	err := setUpViper()
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	err = fromEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setUpViper() error {
	viper.AddConfigPath("./config")
	viper.SetConfigName("main")

	return viper.ReadInConfig()
}

func unmarshal(cfg *Config) error {
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("repo", &cfg.Repo)
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("promo", &cfg.Promo)
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("payment.url", &cfg.Payment.URL)
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("bot.shop", &cfg.Payment.Products)
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("tesseract", &cfg.Tesseract)
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("chat_ai", &cfg.ChatAI)
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("bot", &cfg.Bot)
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("bot.messages.response", &cfg.Bot.Messages.Responses)
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("bot.messages.error", &cfg.Bot.Messages.Errors)
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("bot.keyboard.menu", &cfg.Bot.Keyboard.Menu)
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("bot.keyboard.balance", &cfg.Bot.Keyboard.Balance)
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("bot.keyboard.advices", &cfg.Bot.Keyboard.Advices)
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("bot.commands", &cfg.Bot.Commands)
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("server", &cfg.Server)
	if err != nil {
		return err
	}

	return nil
}

func fromEnv(cfg *Config) error {
	err := viper.BindEnv("TELEGRAM_TOKEN")
	if err != nil {
		return err
	}
	cfg.TelegramToken = viper.GetString("TELEGRAM_TOKEN")

	err = viper.BindEnv("CHAT_GPT_API_KEY")
	if err != nil {
		return err
	}
	cfg.ChatGptApiKey = viper.GetString("CHAT_GPT_API_KEY")

	err = viper.BindEnv("UMONEY_WALLET")
	if err != nil {
		return err
	}
	cfg.Payment.Wallet = viper.GetString("UMONEY_WALLET")

	err = viper.BindEnv("KEY_PAYMENT")
	if err != nil {
		return err
	}
	cfg.Payment.KeyPayment = viper.GetString("KEY_PAYMENT")

	err = viper.BindEnv("IP")
	if err != nil {
		return err
	}
	cfg.Bot.Ip = viper.GetString("IP")

	return nil
}
