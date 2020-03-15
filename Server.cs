using Nancy;

namespace telegram{
    public class Server: NancyModule{
        public Server(){
            Get("/", _ => "ciao mondo!");
        }
    }

}