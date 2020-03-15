using System;
using Telegram.Bot;

namespace telegram
{
    class Program
    {
        static void Main(string[] args)
        {
            var bot = new TelegramBotClient(Secrets.telegram_secret);
            var telegram_client = bot.GetMeAsync().Result;
            Console.WriteLine($"Ciao sono l'utente {telegram_client.Id} e il mio nome è {telegram_client.FirstName}");
            
        }
    }
}
