using System;
using System.Text.Json;

namespace demo
{
    class Program
    {
        static void Main(string[] args)
        {
            // replace to your access_id and secret_key
            var accessId = "";
            var secretKey = "";
            var client = new CoinExClient(accessId, secretKey);

            JsonElement result;

            result = client.MarketList();
            Console.WriteLine(result);

            result = client.MarketDepth("BTCUSDT", "1");
            Console.WriteLine(result);

            result = client.AccountInfo();
            Console.WriteLine(result);
            // buy 0.01 BTC use USDT at price 3000
            result = client.PlaceLimitOrder("BTCUSDT", "buy", "0.01", "3000");
            Console.WriteLine(result);
        }
    }
}
