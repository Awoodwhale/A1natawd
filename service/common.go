package service

import (
	"fmt"
	"go_awd/conf"
	"go_awd/pkg/wlog"
	"gopkg.in/gomail.v2"
	"path"
	"runtime"
	"strings"
)

func Debugln(args ...string) {
	wlog.Logger.Debugln(getFuncName(args...)...)
}

func Infoln(args ...string) {
	wlog.Logger.Infoln(getFuncName(args...)...)
}

func Errorln(args ...string) {
	wlog.Logger.Errorln(getFuncName(args...)...)
}

func getFuncName(args ...string) []any {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return nil
	}
	funcNameArr := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	funcName := funcNameArr[len(funcNameArr)-1]
	fileWithLine := fmt.Sprintf("[%s|%v|%s]", path.Base(file), line, funcName)
	args[0] = fileWithLine + args[0]
	res := make([]any, len(args))
	for i, s := range args {
		res[i] = s
	}
	return res
}

// EmailSender
// @Description: 邮箱发送器
type EmailSender struct {
	EmailHost   string
	EmailPort   int
	AuthCode    string
	SendName    string
	FromEmail   string
	ToEmail     string
	Subject     string
	Content     string
	ContentType string
}

func NewEmailSender(toEmail, subject, content, contentType string) *EmailSender {
	return &EmailSender{
		EmailHost:   conf.SmtpHost,
		EmailPort:   conf.SmtpPort,
		AuthCode:    conf.SmtpToken,
		SendName:    conf.SmtpSendName,
		FromEmail:   conf.SmtpEmail,
		ToEmail:     toEmail,
		Subject:     subject,
		Content:     content,
		ContentType: contentType,
	}
}

func (e *EmailSender) SendEmail() error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(e.FromEmail, e.SendName))
	m.SetHeader("To", e.ToEmail)
	m.SetHeader("Subject", e.Subject)
	m.SetBody(e.ContentType, e.Content)
	return gomail.NewDialer(e.EmailHost, e.EmailPort, e.FromEmail, e.AuthCode).DialAndSend(m)
}

func SendUserRegisterEmail(toEmail, captcha string) error {
	content := fmt.Sprintf(template, "请验证您的 A1natawd 帐户", "感谢您注册 A1natawd 帐户，您的验证码如下，请在 5 分钟内进行验证。如果该验证码不为您本人申请，请无视。", captcha)
	return NewEmailSender(toEmail, "A1natawd用户注册通知", content, "text/html; charset=UTF-8").SendEmail()
}

func SendUserUpdateEmail(toEmail, captcha string) error {
	content := fmt.Sprintf(template, "请验证您的 A1natawd 帐户", "您正在进行修改 A1natawd 帐户邮箱的操作，您的验证码如下，请在 5 分钟内进行验证。如果该验证码不为您本人申请，请无视。", captcha)
	return NewEmailSender(toEmail, "A1natawd修改邮箱通知", content, "text/html; charset=UTF-8").SendEmail()
}

func SendUserUpdatePwd(toEmail, captcha string) error {
	content := fmt.Sprintf(template, "请验证您的 A1natawd 帐户", "您正在进行修改 A1natawd 帐户密码的操作，您的验证码如下，请在 5 分钟内进行验证。如果该验证码不为您本人申请，请无视。", captcha)
	return NewEmailSender(toEmail, "A1natawd修改密码通知", content, "text/html; charset=UTF-8").SendEmail()
}

func SendUserRecoverPwd(toEmail, captcha string) error {
	content := fmt.Sprintf(template, "请验证您的 A1natawd 帐户", "您正在进行找回 A1natawd 帐户密码的操作，您的验证码如下，请在 5 分钟内进行验证。如果该验证码不为您本人申请，请无视。", captcha)
	return NewEmailSender(toEmail, "A1natawd找回密码通知", content, "text/html; charset=UTF-8").SendEmail()
}

func SendUserRecoveredPwd(toEmail, pwd string) error {
	content := fmt.Sprintf(template, "请查收您的 A1natawd 帐户", "您正在进行恢复 A1natawd 帐户密码的操作，您的新密码如下，请在收到本邮箱后进行修改密码的操作，确保账号安全。", pwd)
	return NewEmailSender(toEmail, "A1natawd恢复密码通知", content, "text/html; charset=UTF-8").SendEmail()
}

func SendUserAcceptTeam(toEmail, teamName string) error {
	content := fmt.Sprintf(template, "A1natawd 入队申请通知", "恭喜您，您申请入队的【"+teamName+"】团队，通过了您的入队申请！", teamName+" 欢迎您的加入")
	return NewEmailSender(toEmail, "A1natawd团队通知", content, "text/html; charset=UTF-8").SendEmail()
}

func SendUserRejectTeam(toEmail, teamName string) error {
	content := fmt.Sprintf(template, "A1natawd 入队申请通知", "对不起，您申请入队的【"+teamName+"】团队，拒绝了您的入队申请。", teamName+" 尝试寻找其他团队吧！")
	return NewEmailSender(toEmail, "A1natawd团队通知", content, "text/html; charset=UTF-8").SendEmail()
}

const template = `<table border="0" width="100%%" cellpadding="0" cellspacing="0" bgcolor="#f8f8f8">
    <tbody>
    <tr>
    </tr>
    <tr>
        <td height="30" style="font-size: 30px; line-height: 30px;">&nbsp;</td>
    </tr>
    <tr>
        <td height="30" style="font-size: 30px; line-height: 30px;">
            <div style="text-align: center;margin-bottom: 10px;">
                <img height="50"
                     src="https://pic.imgdb.cn/item/64115baaebf10e5d53c4659a.png">
            </div>
        </td>
    </tr>
    <tr>
        <td align="center">
            <table border="0" align="center" width="590" cellpadding="0" cellspacing="0" bgcolor="#ffffff"
                   style="border-collapse:collapse; mso-table-lspace:0pt; mso-table-rspace:0pt;" class="container590">
                <tbody>
                <tr>
                    <td height="45" style="font-size: 45px; line-height: 45px;">&nbsp;</td>
                </tr>
                <tr>
                    <td align="center"
                        style="color: #222222; font-size: 24px; font-family: 'Neue,Helvetica,PingFang SC,Tahoma,Arial,sans-serif;', Arial, sans-serif; mso-line-height-rule: exactly;"
                        class="cta-header">
                        <div>
                            %s
                        </div>
                    </td>
                </tr>
                <tr>
                    <td height="25" style="font-size: 25px; line-height: 25px;">&nbsp;</td>
                </tr>
                <tr>
                    <td>
                        <table border="0" align="center" width="490" cellpadding="0" cellspacing="0" bgcolor="#FFFFFF"
                               class="container580">
                            <tbody>
                            <tr>
                                <td align="center"
                                    style="color: #adb3ba; font-size: 16px; font-weight: 300; font-family: 'Source Sans Pro', Arial, sans-serif; mso-line-height-rule: exactly; line-height: 24px;">
                                    <div style="line-height: 24px;">
                                        %s
                                    </div>
                                </td>
                            </tr>
                            </tbody>
                        </table>
                    </td>
                </tr>
                <tr>
                    <td height="25" style="font-size: 25px; line-height: 25px;">&nbsp;</td>
                </tr>
                <tr>
                    <td align="center">
                        <div style="display: block; width: 550px; height: 40px; border-style: none !important; border: 0 !important;">
                            <div style="text-align:center; font-size:25px;">%s</div>
                        </div>
                    </td>
                </tr>
                <tr>
                    <td height="45" style="font-size: 45px; line-height: 45px;">&nbsp;</td>
                </tr>
                </tbody>
            </table>
        </td>
    </tr>
    <tr>
        <td height="100" style="font-size: 100px; line-height: 100px;">&nbsp;</td>
    </tr>
    </tbody>
</table>`
