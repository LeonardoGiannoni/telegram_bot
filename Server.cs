using Nancy;
using System.Threading.Tasks;
using Nancy.Extensions;
using Telegram.Bot.Types;

namespace telegram{
    public class Server: NancyModule{
        public Server(){
            Get("/", _ => {
                
                return "ciao mondo!";
            });

            Post("/", _ =>
            {
                var jsonString = this.Request.Body.AsString();
                Message message= await botClient.SendTextMessageAsync(
                channelId: Secrets.telegram_channel,
                text: "Test message to a Bot telegram",
                parseMode: ParseMode.Markdown,
                disableNotification: true,
                replyToMessageId: e.Message.MessageId     
                )
            });
        }
    }
}