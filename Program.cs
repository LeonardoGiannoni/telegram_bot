    using System;
    using Telegram.Bot;
    using Nancy.Hosting.Self;
    using System.Threading.Tasks;

//netsh http add urlacl url="http://+:12345/" user="Everyone"
namespace telegram
{
    class Program
    {
        static int Main(string[] args)
        {
            //--client part
            var bot = new TelegramBotClient(Secrets.telegram_secret);
            var telegram_client = bot.GetMeAsync().Result;
            Console.WriteLine($"Ciao sono l'utente {telegram_client.Id} e il mio nome è {telegram_client.FirstName}");
            
            //--Server part--
            string server_hostname = "localhost";
            int port = 8081;
            var config = new HostConfiguration()//prenota namespace dell'URL per utenti non admin 
            {
                UrlReservations = new UrlReservations() { CreateAutomatically = true }
            };
            var server_telegram  = new NancyHost(config, new Uri($"http://{server_hostname}:{port}"));
            server_telegram.Start();
            Console.WriteLine(String.Format($"Server started @ http://{server_hostname}:{port}"));
            Console.ReadKey();
            return 0;
        }
    }
}
