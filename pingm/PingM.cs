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
        private const string APP_NAME = "pingm";
        private const int MIN_ARG_SIZE = 2;

        static int Main(string[] args)
        {
            // Don't run if nothing supplied from user
            if (args.Length < MIN_ARG_SIZE)
            {
                PrintHelp();
                Environment.Exit(1);
            }

            // TODO Add error handling
            bool isRunning = true;
            
            if (!int.TryParse(args[0], out int timeOut))
            {
                PrintInvalidTimeout();
                Environment.Exit(1);
            }

            timeOut *= 1_000;   // Convert into millis

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

            // Copy nodes out of args and convert to 'NetworkNode' type
            // Filter out items that don't resolve and let user know not valid
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

                    // Try resolve to IP or display if it does not resolve
                    try
                    {
                        dns = Dns.GetHostEntry(potentialNode).AddressList;
                    }
                    catch (SocketException)
                    {
                        PrintNotValidNode(potentialNode);
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


        /// <summary>
        /// 
        /// </summary>
        private static void PrintNoValidNodes()
        {
            Console.WriteLine("No valid hosts supplied.");
        }

        /// <summary>
        /// 
        /// </summary>
        private static void PrintInvalidTimeout()
        {
            Console.WriteLine("Invalid timeout value supplied.");
        }


        /// <summary>
        /// 
        /// </summary>
        /// <param name="potentialNode"></param>
        private static void PrintNotValidNode(string potentialNode)
        {
            Console.WriteLine($"Hostname '{potentialNode}' does not resolve to an IP address.");
        }


        /// <summary>
        /// 
        /// </summary>
        private static void PrintHelp()
        {
            Console.WriteLine($"{APP_NAME} <timeout in seconds> <host1> <host2> <host..> <host10>");
        }


        /// <summary>
        /// 
        /// </summary>
        private static void PrintHeader()
        {
            Console.WriteLine("{0,-20} {1,-15}", "Hostname", "IP Address");
        }


        /// <summary>
        /// 
        /// </summary>
        /// <param name="node"></param>
        /// <param name="timeOut"></param>
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

