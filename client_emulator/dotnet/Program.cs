// #define ACCOUNT_LOGIN
// #define S3_PUT_KEY
#define S3_GET_KEY

using System.Net.Http.Headers;
using Google.Protobuf;
using Protocol;

using var http = new HttpClient();

#region account login
#if ACCOUNT_LOGIN
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
#endif
#endregion

#region s3 put key
#if S3_PUT_KEY
var request = new S3PutKeyReq { Key = "dotnet-test-key", Data = "dotnet-test-data" };
var url = "http://localhost:8084/s3/PutKey";
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
var putKeyResp = S3PutKeyResp.Parser.ParseFrom(Convert.FromBase64String(message!));
Console.WriteLine($"PutKeyResp: Code={putKeyResp.Code}, Message={putKeyResp.Message}");
#endif
#endregion

#region s3 get key
#if S3_GET_KEY
var request = new S3GetKeyReq { Key = "dotnet-test-key" };
var url = "http://localhost:8084/s3/GetKey";
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
var getKeyResp = S3GetKeyResp.Parser.ParseFrom(Convert.FromBase64String(message!));
Console.WriteLine($"GetKeyResp: Code={getKeyResp.Code}, Data={getKeyResp.Data}, Message={getKeyResp.Message}");
#endif
#endregion