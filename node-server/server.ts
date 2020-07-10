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

/**
 * Implements the generateRandomNumber RPC method.
 */
function generateRandomNumber(call, callback) {
  const {maxLength} = call.request;
  const number = Math.floor(Math.random() * Math.floor(maxLength));

  console.info(`generateRandomNumber method received request. Max number Length = ${maxLength}`);
  console.info(`returning value ${number}`);

  callback(null, {number});
}

/**
 * Starts an RPC server that receives requests for the Greeter service at the
 * sample server port
 */
function main() {
  const server = new grpc.Server();
  server.addService(randomProto.Random.service, {generateRandomNumber: generateRandomNumber});
  server.bind("0.0.0.0:50052", grpc.ServerCredentials.createInsecure());
  server.start();
  console.info("nodejs grpc service started on localhost:50052");
}

main();
