package worker

const EmailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }
        :root{
            font-size: 12px;
        }
        body {
            font-family: Arial, sans-serif;
        }
        .header {
            width: 100%%;
            max-width: 600px;
        }
        .table-content {
            width: 100%%;
            max-width: 600px;
            border: 0;
            border-spacing: 0;
            border-collapse: separate;
        }
        .content {
            font-size: 14px;
            padding: 2rem;
            background-color: #e7f3f3;
        }

        .content .button-wrapper {
            width: 100%%;
            text-align: center;
        }

        .content .button {
            display: inline-block;
            color: #f8f6f8;
            padding: 1rem;
            margin: 1rem;
            border-radius: 8px;
            background-color: #36866f;
            text-decoration: none;
            font-weight: bold;
        }

        .footer {
            width: 100%%;
            min-height: 70px;
            background-color: #202729;
            -webkit-mask-image: linear-gradient(200deg, rgba(255, 255, 255, 0) -200%%, rgb(255, 255, 255) 100%%);
            color: #f8f6f8;
            padding: 12px;
            display: flex;
            flex-direction: row;
        }

        .logo img {
            width: 64px;
            height: auto;
        }

        .socials {
            display: flex;
            flex-direction: row;
            justify-content: center;
            align-items: center;
            margin-left: auto;
            gap: 1rem;
        }

        .socials .icons {
            display: flex;
            flex-direction: row;
            font-size: 1rem;
            gap: .5rem;
            align-items: center;
        }

        .socials .icons a {
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        .socials .icons img{
            width: 2rem;
            height: 2rem;
        }

        .footer a {
            color: #c2e0e0;
        }

        .footer a:hover {
            opacity: 0.7;
        }
    </style>
</head>
<body>
<table class="table-content">
    <tr>
        <td>
            <img src="https://yt3.googleusercontent.com/K5AwRDNYCOGODqJDCJnexGDeWz-T139vHalPRaIKDdLR8nmw8wRzk6WkMespBsJH6S4g9Pyr=w2120-fcrop64=1,00005a57ffffa5a8-k-c0xffffffff-no-nd-rj" alt="Header Image" class="header" style="display: block; width: 100%%; max-width: 600px; height: auto;">
        </td>
    </tr>
    <tr>
        <td>
            <div class="content">
				%s
            </div>
        </td>
    </tr>
    <tr>
        <td>
            <div class="footer">
                <div class="logo">
                    <img alt=""src="https://res.cloudinary.com/db8w4lhqr/image/upload/v1683557020/logo-v1-128_vdwqzr.png">
                </div>
                <div class="socials">
                    <div class="icons">
                        <a target="_blank" href="https://discord.gg/wWsXB5gVmY"><img src="https://img.icons8.com/?size=1x&id=2mIgusGquJFz&format=png" alt="discord" />Discord</a> |
                        <a target="_blank" href="https://www.youtube.com/@Talebound"><img src="https://img.icons8.com/?size=1x&id=9a46bTk3awwI&format=png" alt="youtube" />Youtube</a> |
                        <a target="_blank" href="https://github.com/Talebound"><img src="https://img.icons8.com/?size=1x&id=AZOZNnY73haj&format=png" alt="github" />GitHub</a> |
                        <a target="_blank" href="https://ko-fi.com/talebound"><img src="https://img.icons8.com/?size=1x&id=pDSKjMcOYsWo&format=png" alt="ko-fi" />Ko-Fi</a> |
                    </div>
                    <p>&copy; 2023 Talebound</p>
                </div>
            </div>
        </td>
    </tr>
</table>
<!-- rest of your email content -->
</body>
</html>
`
