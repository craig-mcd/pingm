using System;
using System.Net;

namespace pingm
{
    public struct NetworkNode : IEquatable<NetworkNode>
    {
        public string HostName { get; private set; }
        public IPAddress IP { get; private set; }

        public NetworkNode(string hostName, IPAddress ip)
        {
            HostName = hostName;
            IP = ip;
        }


        public override bool Equals(object obj)
        {
            if ((obj == null) || GetType() != obj.GetType())
            {
                return false;
            }
            else
            {
                return Equals((NetworkNode) obj);
            }
        }


        public bool Equals(NetworkNode node)
        {
            return (HostName == node.HostName) && (IP == node.IP);
        }


        public override int GetHashCode()
        {
            return HostName.GetHashCode() ^ IP.GetHashCode();
        }


        public static bool operator ==(NetworkNode left, NetworkNode right)
        {
            return left.Equals(right);
        }


        public static bool operator !=(NetworkNode left, NetworkNode right)
        {
            return !(left == right);
        }


        public override string ToString()
        {
            return $"Hostname{HostName}, IP: {IP}";
        }
    }
}