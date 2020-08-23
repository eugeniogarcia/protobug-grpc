using Grpc.Net.Client;
using Helloworld;
using System;

namespace cliente_grpc
{
    class Program
    {
        static async System.Threading.Tasks.Task Main(string[] args)
        {
            AppContext.SetSwitch("System.Net.Http.SocketsHttpHandler.Http2UnencryptedSupport", true);
            var channel=GrpcChannel.ForAddress("http://localhost:50051");
            var client = new Greeter.GreeterClient(channel);

            var response = await client.SayHelloAsync(new HelloRequest
            {
                Name= "Eugenio" 
            });

            Console.WriteLine("From server: "+response.Message);
        }
    }
}
