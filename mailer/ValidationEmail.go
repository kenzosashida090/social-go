package mailer

var EmailValidationTemplate = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Welcome</title>
</head>
<body style="font-family: Arial, sans-serif; background:#F6FAFB; margin: 0; padding: 0;">
  <table width="100%" cellspacing="0" cellpadding="0" border="0" bgcolor="#F6FAFB">
    <tr>
      <td align="center">
        <table width="600" cellspacing="0" cellpadding="0" border="0">
          <tr>
            <td style="padding:48px 0 30px; text-align:center; font-size:14px; color:#4C83EE;">
              KibSocial
            </td>
          </tr>

          <tr>
            <td style="padding:48px 30px 40px; color:#000000;" bgcolor="#ffffff">

              <div style="font-size:18px; line-height:150%; font-weight:bold; padding-bottom:24px;">
                Welcome, {{.UserName}}!
              </div>

              <div style="font-size:14px; line-height:150%; padding-bottom:10px;">
                Thanks for using KIbSocial as your post of integrity c;
              </div>

              <div style="font-size:14px; line-height:150%; padding-bottom:16px;">
							To get started, please validate your account by clicking the button bellow
              </div>

              <div style="padding-bottom:24px;">
                <a href="{{.ActivateUrl}}"
                   style="display:block; background:#4C83EE; text-decoration:none;
                          padding:10px 0; color:#fff; font-size:14px;
                          text-align:center; font-weight:bold; border-radius:7px;">
                  Next Step
                </a>
              </div>

              <div style="border-top:1px solid #8B949F; width:117px; margin-bottom:16px;"></div>

              <div style="font-size:14px; line-height:170%;">
                Best regards,<br>
                <strong>KibSocial</strong>
              </div>

            </td>
          </tr>

          <tr>
            <td style="padding:24px 0 48px; text-align:center; font-size:11px; color:#8B949F;">
              KIBSocial<br>
                          </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>
`
