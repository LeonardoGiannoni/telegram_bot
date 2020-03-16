using Nancy;
using System.Threading.Tasks;
using Nancy.Extensions;
using Telegram.Bot.Types;

namespace telegram
{
    public class Server : NancyModule
    {
        public Server()
        {
            Get("/", _ =>
            {

                return "ciao mondo!";
            });

            Post("/", _ =>
            {
                var jsonString = this.Request.Body.AsString();
                var t = Program.bot.SendTextMessageAsync(
                    Secrets.telegram_channel,
                    "Test message to a Bot telegram"
                );
                return null;
            });
        }
    }
}