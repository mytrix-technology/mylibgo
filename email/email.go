package email

//type NotifService interface {
//	Send(ctx context.Context, payload byte) error
//}
//
//type EmailClient struct {
//	cfg *config.App
//}
//
//func NewEmailClient(cfg *config.App) NotifService {
//	return EmailClient{
//		//cfg: cfg,
//	}
//}
//
//func (e EmailClient) Send(ctx context.Context, payload byte) error {
//	b, err := json.Marshal(payload)
//	if err != nil {
//		c.Logger(ctx).Error(fmt.Sprintf("Marshaling %v", err.Error()))
//		return err
//	}
//
//	currTime := time.Now()
//	formattedCurrTime := currTime.Format(model.LongDateFormat3)
//	correlationId := uuid.NewV4().String()
//	stanId := uuid.NewV4().String()
//
//	data := []byte(string(b))
//
//	req, err := http.NewRequest("POST", e.cfg.EMAILRESTHOST+"/2.0.0/send", bytes.NewBuffer(data))
//
//	headers := http.Header{
//		"Content-Type":             []string{"application/json"},
//		"X-Channel-Id":             []string{"6024"},
//		"X-Node":                   []string{"BTPNS"},
//		"X-Correlation-Id":         []string{correlationId},
//		"X-Stan-Id":                []string{stanId},
//		"X-Transmission-Date-Time": []string{formattedCurrTime},
//		"X-Terminal-Id":            []string{"358525071384733"},
//		"X-Terminal-Name":          []string{"bepcap01"},
//		"X-Acq-Id":                 []string{"547"},
//		"X-orgUnit-Id":             []string{"547"},
//		"X-API-Key":                []string{e.cfg.XAPIKEY},
//	}
//	req.Header = headers
//	if strings.Contains(e.cfg.EMAILRESTHOST, "nww") {
//		requestDump, err := httputil.DumpRequest(req, true)
//		if err != nil {
//			c.Logger(ctx).Error(fmt.Sprintf("Fail dumping request %v", err.Error()))
//			return err
//		}
//		c.Logger(ctx).Info("Email Notif Request", zap.Any("request: ", string(requestDump)))
//	} else {
//		body, err := common.DumpBodyRequest(req)
//		if err != nil {
//			c.Logger(ctx).Error(fmt.Sprintf("Fail dumping body request %v", err.Error()))
//			return err
//		}
//		c.Logger(ctx).Info("Email Notif Request",
//			zap.Any("method", req.Method),
//			zap.Any("url", req.URL),
//			zap.Any("body", string(body)),
//			zap.Any("X-Correlation-Id", correlationId),
//			zap.Any("X-Stan-Id", stanId),
//		)
//	}
//	res, err := http.DefaultClient.Do(req)
//	if err != nil {
//		c.Logger(ctx).Error(fmt.Sprintf("Fail send email %v", err.Error()))
//		return err
//	}
//	responseDump, err := httputil.DumpResponse(res, true)
//	c.Logger(ctx).Info("Email Notif Response", zap.Any("response: ", string(responseDump)))
//	defer res.Body.Close()
//	if err != nil || res.StatusCode != 200 {
//		bodyBytes, err := ioutil.ReadAll(res.Body)
//		if err != nil {
//			c.Logger(ctx).Error(fmt.Sprintf("Success send but not 200 error Read %v", err.Error()))
//			return err
//		}
//		c.Logger(ctx).Error(fmt.Sprintf("Success send but not 200 %v", string(bodyBytes)))
//		return err
//	}
//	return err
//}
