using Nancy;
using System.Threading.Tasks;
using Nancy.Extensions;

namespace telegram{
    public class Server: NancyModule{
        public Server(){
            Get("/", _ => "ciao mondo!");

            Post("/", _ =>
            {
                var jsonString = this.Request.Body.AsString();
                
            });
        }
    }
}