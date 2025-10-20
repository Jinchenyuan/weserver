using System.Net.Http.Headers;
using Google.Protobuf;
using Protocol;

using var http = new HttpClient();

#region account login
var request = new AccountLoginReq { Username = "linux-user", Password = "password" };
var url = "http://10.4.11.146:8083/account/login";
var body = new
{
    Message = Convert.ToBase64String(request.ToByteArray())
};
var json = System.Text.Json.JsonSerializer.Serialize(body);
var content = new StringContent(json);
content.Headers.ContentType = new MediaTypeHeaderValue("application/json");
var response = await http.PostAsync(url, content);

// read response
var responseBody = await response.Content.ReadAsStringAsync();
var responseJson = System.Text.Json.JsonDocument.Parse(responseBody);
var message = responseJson.RootElement.GetProperty("message").GetString();
Console.WriteLine(message);
var loginResp = AccountLoginResp.Parser.ParseFrom(Convert.FromBase64String(message!));
Console.WriteLine($"LoginResp: Code={loginResp.Code}, Token={loginResp.Token}, Message={loginResp.Message}");
#endregion