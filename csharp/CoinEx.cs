using System;
using System.Text;
using System.Collections.Generic;
using System.Net.Http;
using System.Text.Json;
using System.Linq;
using System.Security.Cryptography;
using System.Net.Http.Headers;

namespace demo
{
    class CoinExClient
    {
        private string baseUrl;
        private string accessId;
        private string secretKey;
        private HttpClient client;

        public CoinExClient(string accessId, string secretKey)
        {
            this.baseUrl = "https://api.coinex.com";
            this.accessId = accessId;
            this.secretKey = secretKey;
            this.client = new HttpClient();
        }

        private string Sign(Dictionary<string, object> args)
        {
            // add access_id and tonce to args
            args["access_id"] = this.accessId;
            args["tonce"] = (long)(DateTime.UtcNow - new DateTime(1970, 1, 1)).TotalMilliseconds;
            // sort args by name and add secret_key to the last
            var last_args = new[] { new KeyValuePair<string, object>("secret_key", this.secretKey) };
            var sortedArgs = args.OrderBy(p => p.Key).Concat(last_args);
            var body = string.Join("&", sortedArgs.Select(p => $"{p.Key}={p.Value}"));
            //calc md5
            var md5 = MD5.Create().ComputeHash(Encoding.UTF8.GetBytes(body));
            return BitConverter.ToString(md5).Replace("-", string.Empty).ToUpper();
        }

        private JsonElement Post(string path, Dictionary<string, object> args, string signature = null)
        {
            var json = JsonSerializer.Serialize(args);
            var content = new StringContent(json);
            content.Headers.ContentType = new MediaTypeHeaderValue("application/json");
            var req = new HttpRequestMessage(HttpMethod.Post, this.baseUrl + path);
            req.Content = content;
            if (signature != null)
            { 
                req.Headers.Add("authorization", signature);
            }
            var res = this.client.SendAsync(req).Result;
            var result = res.Content.ReadAsStringAsync().Result;
            return JsonSerializer.Deserialize<JsonElement>(result);
        }

        private JsonElement Get(string path, Dictionary<string, object> args = null, string signature = null)
        {
            var url = this.baseUrl + path;
            if (args != null)
            {
                var param = string.Join("&", args.Select(p => $"{p.Key}={p.Value}"));
                url = url + "?" + param;
            }
            string result;
            if (signature != null)
            {
                var req = new HttpRequestMessage(HttpMethod.Get, url);
                req.Headers.Add("authorization", signature);
                var res = client.SendAsync(req).Result;
                result = res.Content.ReadAsStringAsync().Result;
            }
            else
            {
                result = this.client.GetStringAsync(url).Result;
            }
            return JsonSerializer.Deserialize<JsonElement>(result);
        }

        public JsonElement MarketList()
        {
            return this.Get("/v1/market/list");
        }

        public JsonElement MarketDepth(string market, string merge)
        {
            var args = new Dictionary<string, object>
            {
                ["market"] = market,
                ["merge"] = merge
            };
            return this.Get("/v1/market/depth", args);
        }

        public JsonElement AccountInfo()
        {
            var args = new Dictionary<string, object>();
            var signature = this.Sign(args);
            return this.Get("/v1/balance/info", args, signature);
        }

        public JsonElement PlaceLimitOrder(string market, string type, string amount, string price)
        {
            var args = new Dictionary<string, object>
            {
                ["market"] = market,
                ["type"] = type,
                ["amount"] = amount,
                ["price"] = price
            };
            var signature = this.Sign(args);
            return this.Post("/v1/order/limit", args, signature);
        }
    }
}
