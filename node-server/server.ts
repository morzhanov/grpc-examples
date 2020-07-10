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
 * Implements the streamNumbers RPC method.
 */
function streamNumbers(call) {
  const {maxLength} = call.request;
  const number = Math.floor(Math.random() * Math.floor(maxLength));

  console.info(`streamNumbers method received request. Max number Length = ${maxLength}`);
  console.info(`streaming data...`);

  for (let i = 0; i < 2; ++i) {
    call.write({number});
  }
  call.end();
}

/**
 * Implements the logStreamOfRandomNumbers RPC method.
 */
function logStreamOfRandomNumbers(call, callback) {
  call.on("data", function (number) {
    console.info(`logStreamOfRandomNumbers got number: ${number}`);
  });
  call.on("end", function () {
    console.info(`logStreamOfRandomNumbers finished`);
    callback(null, {message: "Successfully logged all numbers"});
  });
}

/**
 * Implements the bidirectionalStream RPC method.
 */
function bidirectionalStream(call) {
  call.on("data", function ({number}) {
    const resNumber = Math.random() * 100;

    console.info(`bidirectionalStream got number: ${number}`);
    console.info(`bidirectionalStream pushed number: ${resNumber}`);

    call.write({number: resNumber});
  });
  call.on("end", function () {
    console.info(`bidirectionalStream finished`);
    call.end();
  });
}

/**
 * Starts an RPC server that receives requests for the Greeter service at the
 * sample server port
 */
function main() {
  const server = new grpc.Server();
  server.addService(randomProto.Random.service, {
    generateRandomNumber,
    streamNumbers,
    logStreamOfRandomNumbers,
    bidirectionalStream,
  });

  server.bind("0.0.0.0:50052", grpc.ServerCredentials.createInsecure());
  server.start();
  console.info("nodejs grpc service started on localhost:50052");
}

main();
