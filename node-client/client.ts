import * as grpc from "grpc";
import * as protoLoader from "@grpc/proto-loader";
import {resolve} from "path";

const PROTO_PATH = resolve(__dirname, "../proto/random.proto");

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});
const randomProto: any = grpc.loadPackageDefinition(packageDefinition).random;

function getRundomNumber(client) {
  client.generateRandomNumber({maxLength: process.argv[2]}, function (err, response) {
    if (err) {
      console.error(err);
    }
    console.log("Random number: ", response.number);
  });
}

function main() {
  const address = process.argv[3] ? `localhost:${process.argv[3]}` : "localhost:50052";
  console.info(`Connecting to server: ${address}`);

  var client = new randomProto.Random(address, grpc.credentials.createInsecure());

  setInterval(() => getRundomNumber(client), 1000);
}

main();
