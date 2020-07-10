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

async function getRundomNumber(client) {
  return new Promise((resolve, reject) => {
    client.generateRandomNumber({maxLength: process.argv[2]}, function (err, response) {
      if (err) {
        console.error(err);
        reject(err);
      }
      console.log("Random number: ", response.number);
      resolve();
    });
  });
}

async function streamNumbers(client) {
  return new Promise((resolve, reject) => {
    var call = client.streamNumbers({maxLength: process.argv[2]});

    call.on("data", function ({number}) {
      console.log(`streamNumbers got number: ${number}`);
    });
    call.on("end", function () {
      console.log(`streamNumbers successfully finished`);
      resolve();
    });
    call.on("error", function (e) {
      console.log(`streamNumbers error ${e}`);
      reject(e);
    });
    call.on("status", function (status) {
      console.log(`streamNumbers status ${status}`);
    });
  });
}

async function logStreamOfRandomNumbers(client) {
  return new Promise(async (resolve, reject) => {
    var call = client.logStreamOfRandomNumbers(function (error, {message}) {
      if (error) {
        console.log(`logStreamOfRandomNumbers error ${error}`);
        reject(error);
      }
      console.log(`logStreamOfRandomNumbers successfully finished`);
      resolve(message);
    });

    for (let i = 0; i < 2; ++i) {
      const number = Math.floor(Math.random() * Math.floor(200));
      await call.write({number});
    }

    call.end();
  });
}

async function bidirectionalStream(client) {
  return new Promise(async (resolve, reject) => {
    var call = client.bidirectionalStream(function (error, {message}) {
      if (error) {
        console.log(`bidirectionalStream error ${error}`);
        reject(error);
      }
      console.log(`bidirectionalStream successfully finished`);
      resolve(message);
    });

    call.on("data", function ({number}) {
      console.log(`bidirectionalStream got number: ${number}`);
    });
    call.on("end", function () {
      console.log(`bidirectionalStream successfully finished`);
      resolve();
    });
    call.on("error", function (e) {
      console.log(`bidirectionalStream error ${e}`);
      reject(e);
    });
    call.on("status", function (status) {
      console.log(`bidirectionalStream status:`);
      console.dir(status);
    });

    for (let i = 0; i < 2; ++i) {
      const number = Math.floor(Math.random() * Math.floor(200));
      await call.write({number});
    }

    call.end();
  });
}

function main() {
  const address = process.argv[3] ? `localhost:${process.argv[3]}` : "localhost:50052";
  console.info(`Connecting to server: ${address}`);

  var client = new randomProto.Random(address, grpc.credentials.createInsecure());

  setInterval(async () => {
    await getRundomNumber(client);
    await streamNumbers(client);
    await logStreamOfRandomNumbers(client);
    await bidirectionalStream(client);
  }, 1000);
}

main();
