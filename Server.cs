using Nancy;

namespace telegram {

    class Server : NancyModule {
        public Server() {
            Get("/") = async (_, token) => 
            {
                return "Hello World";
            };
        }
    }
}