using System;
using System.Net;

namespace pingm
{
    public struct NetworkNode: IEquatable<NetworkNode>
    {
        public string HostName { get; private set; }
        public IPAddress IP { get; private set; }

        public NetworkNode(string hostName, IPAddress ip)
        {
            HostName = hostName;
            IP = ip;
        }

        public bool Equals(NetworkNode node)
        {
            return true;
        }

        public override bool Equals(object obj)
        {
            // throw new System.NotImplementedException();
            return true;
        }

        public override int GetHashCode()
        {
            // throw new System.NotImplementedException();
            return 42;
        }

        public static bool operator ==(NetworkNode left, NetworkNode right)
        {
            return left.Equals(right);
        }

        public static bool operator !=(NetworkNode left, NetworkNode right)
        {
            return !(left == right);
        }
    }
}