package core

import "github.com/tokopedia/sweep-log/core/enum"

const (
	NOTIFY_SUCCESS = 1
	VALIDATE_USE   = 2

	MODE_ONE_BY_ONE = 1
	MODE_ALL_IN_ONE = 2
)

type Filter struct {
	GrepType enum.GrepType
	Value    string
}

type (
	UseCodeInput struct {
		Data UseCodeData `json:"data"`
	}

	UseCodeData struct {
		Codes              []string               `json:"codes"`
		CurrentApplyCode   UseCodeApplyInfo       `json:"current_apply_code"`
		PromoCodeSessionID int64                  `json:"promo_code_id"`
		UserData           UseCodeUserData        `json:"user_data"`
		PaymentInfo        UseCodePaymentInfo     `json:"payment_info"`
		MetaData           map[string]interface{} `json:"meta_data"`
		GrandTotal         float64                `json:"grand_total"`
		PaymentAmount      float64                `json:"payment_amount"`

		// service id and secret key
		ServiceID int64  `json:"service_id"`
		SecretKey string `json:"secret_key"`

		// promocode related
		IsFirstStep bool `json:"is_first_step"`
		Book        bool `json:"book"`
		IsSuggested bool `json:"is_suggested"`
		IsAppRoot   bool `json:"is_app_root"`
		SkipApply   bool `json:"skip_apply"`

		// rule engine related
		Campaign []UseCodeCampaignData `json:"campaign_data"`

		// after payment
		PaymentID                int64 `json:"payment_id"`
		FinishedOrderID          int64 `json:"finished_order_id"`
		IsSkipCheckNotifySuccess bool  `json:"is_skip_check_notify_success"`

		IsSkipFraudCheck bool      `json:"is_skip_fraud_check"`
		FraudInfo        FraudInfo `json:"fraud_info"`

		// others
		Language string `json:"language"`
		State    string `json:"state"`

		//data recovery
		DataRecoveryMode bool `json:"data_recovery_mode"`
	}

	UseCodeUserData struct {
		UserID                int64                  `json:"user_id"`
		Name                  string                 `json:"name"`
		Email                 string                 `json:"email"`
		Msisdn                string                 `json:"msisdn"`
		MsisdnVerified        bool                   `json:"msisdn_verified"`
		IsQcAccount           bool                   `json:"is_qc_acc"`
		AppVersion            string                 `json:"app_version"`
		UserAgent             string                 `json:"user_agent"`
		IPAddress             string                 `json:"ip_address"`
		AdsID                 string                 `json:"advertisement_id"`
		DeviceType            string                 `json:"device_type"`
		DeviceID              string                 `json:"device_id"`
		AddressDetail         map[string]interface{} `json:"address_detail"`
		UserTransactionDetail map[string]interface{} `json:"user_transaction_detail"`
	}

	UseCodePaymentInfo struct {
		ScroogeGatewayID      int64   `json:"scrooge_gateway_id"`
		ScroogeGatewayCode    string  `json:"scrooge_gateway_code"`
		ScroogeGatewayValue   string  `json:"scrooge_gateway_value"`
		CreditCardNumber      string  `json:"credit_card_number"`
		CreditCardExpiryMonth int     `json:"cc_exp_month"`
		CreditCardExpiryYear  int     `json:"cc_exp_year"`
		CreditCardHash        string  `json:"cc_hash"`
		OvoCashAmount         float64 `json:"ovo_cash_amount"`
		OvoPointsAmount       float64 `json:"ovo_points_amount"`
	}

	UseCodeCampaignData struct {
		Code        string                 `json:"code"`
		PromoCodeID int64                  `json:"promo_code_id"`
		PromoID     int64                  `json:"promo_id"`
		RuleIDs     []int64                `json:"rule_ids"`
		DoGaladriel []string               `json:"do_galadriel"`
		FlowState   int                    `json:"flow_state"`
		MetaData    map[string]interface{} `json:"meta_data"`
	}

	PromoStackAbuse struct {
		PromoID       int64 `json:"promo_id"`
		IsPromoAbuser int   `json:"is_promo_abuser"`
	}

	UseCodeServiceData struct {
		CategoryCode int    `json:"category_code"`
		ProductCode  string `json:"product_code"`
	}

	UseCodeApplyInfo struct {
		Code string `json:"code"`
		Type string `json:"type"`
	}
)

type FraudInfo struct {
	DropshipAsBuyer int               `json:"dropship_as_buyer"`
	IsPromoAbuser   int               `json:"is_promo_abuser"`
	IsInvalid       int               `json:"is_invalid"`
	Status          interface{}       `json:"status"`
	PromoStackAbuse []PromoStackAbuse `json:"promo_stack_abuse"`
	Source          string            `json:"source"`
}
