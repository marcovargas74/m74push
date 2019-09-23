<?php
// Source: https://github.com/Hitman666/ios-voip-push

// Put your device token here (without spaces):
$deviceToken = '42c4741429884981d5e949cfb44e571e9fbfecece823565fca9ffb410f558406';
// Put your private key's passphrase here:
$passphrase = '';
// Put your alert message here:
$message = 'Recebendo Chamada Teste';
// Put the full path to your .pem file
$pemFile = 'mobiliti2in1.pem';

////////////////////////////////////////////////////////////////////////////////
$ctx = stream_context_create();
stream_context_set_option($ctx, 'ssl', 'local_cert', $pemFile);
stream_context_set_option($ctx, 'ssl', 'passphrase', $passphrase);
// Open a connection to the APNS server
$fp = stream_socket_client(
    'ssl://gateway.sandbox.push.apple.com:2195', $err,
    $errstr, 60, STREAM_CLIENT_CONNECT|STREAM_CLIENT_PERSISTENT, $ctx);
if (!$fp)
    exit("Failed to connect: $err $errstr" . PHP_EOL);
echo 'Connected to APNS' . PHP_EOL;
// Create the payload body
$body['aps'] = array(
    'alert' => $message,
    'sound' => 'default'
);
// Encode the payload as JSON
$payload = json_encode($body);
// Build the binary notification
$msg = chr(0) . pack('n', 32) . pack('H*', $deviceToken) . pack('n', strlen($payload)) . $payload;
// Send it to the server
$result = fwrite($fp, $msg, strlen($msg));
if (!$result)
    echo 'Message not delivered' . PHP_EOL;
else
    echo 'Message successfully delivered' . PHP_EOL;
// Close the connection to the server
fclose($fp);
