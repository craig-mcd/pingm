using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using System.Net;
using System.Net.NetworkInformation;
using System.Net.Sockets;
using System.Threading;
using System.Text;

// TODO CLI parsing
// TODO Option for log file

namespace pingm
{
    class PingM
    {
        static int Main(string[] args)
        {
            const int MIN_ARG_SIZE = 2;

            // Don't run if nothing supplied from user
            if (args.Length < MIN_ARG_SIZE)
            {
                PrintHelp();
                Environment.Exit(1);
            }

            // TODO Add error handling
            bool isRunning = true;
            int timeOut = int.Parse(args[0]) * 1_000; // Convert from seconds to millis
            var nodes = new List<NetworkNode>();

            // Event handler for CTRL-C
            Console.CancelKeyPress += delegate(object sender, ConsoleCancelEventArgs args)
            {
                Console.ForegroundColor = ConsoleColor.Black;
                Console.BackgroundColor = ConsoleColor.Yellow;
                Console.WriteLine("Finishing...");
                Console.ResetColor();
                args.Cancel = true;
                isRunning = false;
            };

            // Copy nodes out of args and convert to Node type
            // Filter out items that don't resolve
            for (int i = 1; i < args.Length; i++)
            {
                string potentialNode = args[i];
                NetworkNode node;
                bool isIp = IPAddress.TryParse(potentialNode, out var ip);

                if (isIp)
                {
                    node = new NetworkNode("", ip);
                }
                else
                {
                    // I don't want to use IPHostEntry as it doesn't end up using the supplied hostname
                    // Just use it to extract the first IP returned from DNS result set

                    IPAddress[] dns;

                    // Try resolve to IP or display it will be ignored
                    try
                    {
                        dns = Dns.GetHostEntry(potentialNode).AddressList;
                    }
                    catch (SocketException)
                    {
                        PrintNotValid(potentialNode);
                        continue;
                    }

                    // Often multiple IPs returned, just take the first
                    // TODO Look at using all IP's returned, maybe as a CLI option
                    node = new NetworkNode(potentialNode, dns[0]);
                }

                nodes.Add(node);
            }

            // Exit if no valid hosts/nodes
            if (nodes.Count == 0)
            {
                PrintNoValidNodes();
                Environment.Exit(1);
            }

            Console.CursorVisible = false;

            while (isRunning)
            {
                Console.ForegroundColor = ConsoleColor.White;
                Console.BackgroundColor = ConsoleColor.Blue;

                Console.WriteLine(DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss"));
                Console.ResetColor();
                PrintHeader();

                foreach (var node in nodes)
                {
                    var task = new Task(() => ProcessNode(node, timeOut));
                    task.Start();
                }

                // TODO Check what is the best/optimal value to add to the sleep millis
                Thread.Sleep(timeOut + 100);
                Console.Write("\n\n");
            }

            Console.ResetColor();
            Console.CursorVisible = true;
            return 0;
        }


        private static void PrintNoValidNodes()
        {
            Console.WriteLine("No valid hosts supplied.");
        }


        private static void PrintNotValid(string potentialNode)
        {
            Console.WriteLine($"Hostname '{potentialNode}' does not resolve to an IP address.");
        }


        private static void PrintHelp()
        {
            const string help = "pingm <timeout in seconds> <host1> <host2> <host..> <host10>";
            Console.WriteLine(help);
        }


        private static void PrintHeader()
        {
            Console.WriteLine("{0,-20} {1,-15}", "Hostname", "IP Address");
        }


        private static void ProcessNode(NetworkNode node, int timeOut)
        {
            var sb = new StringBuilder();
            sb.Append($"{node.HostName,-20} {node.IP, -15} ");

            using Ping ping = new Ping();
            try
            {
                PingReply reply = ping.Send(node.IP, timeOut);

                if (reply?.Status == IPStatus.Success)
                {
                    sb.Append($"{reply.RoundtripTime}ms");
                }
                else
                {
                    sb.Append($"{reply.Status,-5}");
                }
            }
            catch (PingException e)
            {
                // TODO Better error message handling
                // TODO Better way to handle 'Win32Exception'
                sb.Append(e.GetBaseException().GetType());
            }

            Console.WriteLine(sb);
        }
    }
}
