using System.Net;

namespace pingm
{
    public struct NetworkNode
    {
        public string HostName { get; private set; }
        public IPAddress IP { get; private set; }

        public NetworkNode(string hostName, IPAddress ip)
        {
            HostName = hostName;
            IP = ip;
        }
    }
}