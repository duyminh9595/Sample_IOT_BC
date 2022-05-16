"use strict";

// SDK Library to asset with writing the logic
const { Contract } = require("fabric-contract-api");

class IOTContract extends Contract {
  constructor() {
    super("IOTContract");
    this.TxId = "";
  }

  async beforeTransaction(ctx) {
    // default implementation is do nothing
    this.TxId = ctx.stub.getTxID();
    console.log(`we can do some logging for ${this.TxId}  and many more !!`);
  }

  // register sensor for each farm
  async registrationsensor(ctx, namesensor, type, info) {
    const mspid = await ctx.clientIdentity.getMSPID();
    const sensor = {
      namesensor,
      type,
      info,
      docType: "Sensor",
      mspid,
    };
    await ctx.stub.putState(this.TxId, Buffer.from(JSON.stringify(sensor)));
    console.info("============= END : Create Sensor Success ===========");
  }
  //all sensor in chain
  async allsensorinchain(ctx) {
    let queryString = {};
    queryString.selector = {};
    queryString.selector.docType = "Sensor";
    let queryResults = await this.getQueryResultForQueryString(
      ctx.stub,
      JSON.stringify(queryString)
    );
    return queryResults; //shim.success(queryResults);
  }
  async allsensorbelongtofarm(ctx) {
    const mspid = await ctx.clientIdentity.getMSPID();
    let queryString = {};
    queryString.selector = {};
    queryString.selector.docType = "Sensor";
    queryString.selector.mspid = mspid;
    let queryResults = await this.getQueryResultForQueryString(
      ctx.stub,
      JSON.stringify(queryString)
    );
    return queryResults; //shim.success(queryResults);
  }
  async sendinfofromsensortochain(ctx, sensorid, temp, humd, tds, ph, eco2, tvoc) {
    const deviceAsBytes = await ctx.stub.getState(sensorid);
    if (!deviceAsBytes || deviceAsBytes.length === 0) {
      throw new Error(`${deviceAsBytes} does not exist`);
    }
    const mspid = await ctx.clientIdentity.getMSPID();
    let _keyHelper = new Date();
    let _keyYearAsString = _keyHelper.getFullYear().toString();
    let _keyMonthAsString = _keyHelper.getMonth().toString();
    let _keyDateAsString = _keyHelper.getDate().toString();
    const dataDevice = {
      docType: "DataDevice",
      mspid, txId: this.TxId,datecreated: _keyHelper, sensoridtemp,humd,tds,ph, eco2,tvoc,
    };
    try {
      let indexName = "sensorid-year-month-date-txid";
      let indexKey = await ctx.stub.createCompositeKey(indexName, [
        sensorid,_keyYearAsString,_keyMonthAsString,_keyDateAsString,this.TxId]);
      await ctx.stub.putState(
        indexKey,
        Buffer.from(JSON.stringify(dataDevice))
      );
      return {
        key:sensorid +"-" +_keyYearAsString +"-" +_keyMonthAsString +"-" +_keyDateAsString + "-" +this.TxId,
      };
    } catch (e) {
      throw new Error(`The tx ${this.TxId} can not be stored: ${e}`);
    }
  }

  async getdata(ctx) {
    const args = ctx.stub.getArgs();
    const keyValues = args[1].split("-");
    let keys = [];
    keyValues.forEach((element) => keys.push(element));
    let resultsIterator = await ctx.stub.getStateByPartialCompositeKey(
      "sensorid-year-month-date-txid",
      keys
    );
    const allResults = [];
    while (true) {
      const res = await resultsIterator.next();

      if (res.value) {
        allResults.push(res.value.value.toString("utf8"));
      }
      if (res.done) {
        await resultsIterator.close();
        return allResults;
      }
    }
  }
  async getdataByTimeRange(ctx) {
    // we use the args option
    const args = ctx.stub.getArgs();

    // break condition
    if (args.length !== 3) {
      return JSON.stringify({ error: true });
    }

    // we collect our result
    let allResults = [];

    // compose the selector
    let queryString = {};
    queryString.selector = {};
    queryString.selector.datecreated = {
      $gt: args[1],
      $lt: args[2],
    };
    queryString.sort = [{ datecreated: "asc" }];

    //console.log(queryString)
    // --------------------

    // do the query
    let resultsIterator = await ctx.stub.getQueryResult(
      JSON.stringify(queryString)
    );

    // loop over the results and create the allResults array
    let result = await resultsIterator.next();
    while (!result.done) {
      const strValue = Buffer.from(result.value.value.toString()).toString(
        "utf8"
      );
      let record;
      try {
        record = JSON.parse(strValue);
      } catch (err) {
        console.log(err);
        record = strValue;
      }
      allResults.push({ Key: result.value.key, Record: record });
      result = await resultsIterator.next();
    }

    // return the finale result
    return JSON.stringify(allResults);
  }
  async getalldatasensorbelongtofarm(ctx, sensorid) {
    const mspid = await ctx.clientIdentity.getMSPID();
    let queryString = {};
    queryString.selector = {};
    queryString.selector.docType = "DataDevice";
    queryString.selector.mspid = mspid;
    queryString.selector.sensorid = sensorid;
    let queryResults = await this.getQueryResultForQueryString(
      ctx.stub,
      JSON.stringify(queryString)
    );
    return queryResults; //shim.success(queryResults);
  }

  async storeCs(ctx, revenue, datecreated, cstype) {
    // calc our values
    let _commission = 0;
    let _revenue = parseFloat(revenue);
    let _datecreated = datecreated;

    if (cstype === "reco") {
      // 1 %
      _commission = (_revenue / 100) * 1;
    } else if (cstype === "reve") {
      // 10 %
      _commission = (_revenue / 100) * 10;
    }

    // compose our model
    let model = {
      revenue: _revenue,
      commission: _commission,
      datecreated: _datecreated,
      cstype: cstype,
      txId: this.TxId,
    };

    try {
      // store the composite key with a the value
      let indexName = "year~month~txid";

      let _keyHelper = new Date(datecreated);
      let _keyYearAsString = _keyHelper.getFullYear().toString();
      let _keyMonthAsString = _keyHelper.getMonth().toString();

      let yearMonthIndexKey = await ctx.stub.createCompositeKey(indexName, [
        _keyYearAsString,
        _keyMonthAsString,
        this.TxId,
      ]);

      //console.info(yearMonthIndexKey, _keyYearAsString, _keyMonthAsString, this.TxId);

      // store the new state
      await ctx.stub.putState(
        yearMonthIndexKey,
        Buffer.from(JSON.stringify(model))
      );

      // compose the return values
      return {
        key: _keyYearAsString + "~" + _keyMonthAsString + "~" + this.TxId,
      };
    } catch (e) {
      throw new Error(`The tx ${this.TxId} can not be stored: ${e}`);
    }
  }

  async getdatasensorbelongtofarmonmonthyear(
    ctx,
    sensorid,
    _keyYearAsString,
    _keyMonthAsString
  ) {
    const mspid = await ctx.clientIdentity.getMSPID();
    let queryString = {};
    queryString.selector = {};
    queryString.selector.docType = "DataDevice";
    queryString.selector.mspid = mspid;
    queryString.selector.sensorid = sensorid;
    queryString.selector._keyYearAsString = _keyYearAsString;
    queryString.selector._keyMonthAsString = _keyMonthAsString;
    let queryResults = await this.getQueryResultForQueryString(
      ctx.stub,
      JSON.stringify(queryString)
    );
    return queryResults; //shim.success(queryResults);
  }

  async getdatasensorbelongtofarmondatemonthyear(
    ctx,
    sensorid,
    _keyYearAsString,
    _keyMonthAsString,
    _keyDateAsString
  ) {
    const mspid = await ctx.clientIdentity.getMSPID();
    let queryString = {};
    queryString.selector = {};
    queryString.selector.docType = "DataDevice";
    queryString.selector.mspid = mspid;
    queryString.selector.sensorid = sensorid;
    queryString.selector._keyYearAsString = _keyYearAsString;
    queryString.selector._keyMonthAsString = _keyMonthAsString;
    queryString.selector._keyDateAsString = _keyDateAsString;
    let queryResults = await this.getQueryResultForQueryString(
      ctx.stub,
      JSON.stringify(queryString)
    );
    return queryResults; //shim.success(queryResults);
  }
  async getCsByTimeRange(ctx) {
    // we use the args option
    const args = ctx.stub.getArgs();

    // break condition
    if (args.length !== 3) {
      return JSON.stringify({ error: true });
    }

    // we collect our result
    let allResults = [];

    // compose the selector
    let queryString = {};
    queryString.selector = {};
    queryString.selector.revenueTs = {
      $gt: args[1],
      $lt: args[2],
    };
    queryString.sort = [{ revenueTs: "asc" }];

    //console.log(queryString)
    // --------------------

    // do the query
    let resultsIterator = await ctx.stub.getQueryResult(
      JSON.stringify(queryString)
    );

    // loop over the results and create the allResults array
    let result = await resultsIterator.next();
    while (!result.done) {
      const strValue = Buffer.from(result.value.value.toString()).toString(
        "utf8"
      );
      let record;
      try {
        record = JSON.parse(strValue);
      } catch (err) {
        console.log(err);
        record = strValue;
      }
      allResults.push({ Key: result.value.key, Record: record });
      result = await resultsIterator.next();
    }

    // return the finale result
    return JSON.stringify(allResults);
  }
  async getQueryResultForQueryString(stub, queryString) {
    console.info("- getQueryResultForQueryString queryString:\n" + queryString);
    let resultsIterator = await stub.getQueryResult(queryString);

    let results = await this.getAllResults(resultsIterator, false);

    //return Buffer.from(JSON.stringify(results));
    return results;
  }
  async getAllResults(iterator, isHistory) {
    let allResults = [];
    while (true) {
      let res = await iterator.next();

      if (res.value && res.value.value.toString()) {
        let jsonRes = {};
        console.log(res.value.value.toString("utf8"));

        if (isHistory && isHistory === true) {
          jsonRes.TxId = res.value.tx_id;
          jsonRes.Timestamp = res.value.timestamp;
          jsonRes.IsDelete = res.value.is_delete.toString();
          try {
            jsonRes.Value = JSON.parse(res.value.value.toString("utf8"));
          } catch (err) {
            console.log(err);
            jsonRes.Value = res.value.value.toString("utf8");
          }
        } else {
          jsonRes.Key = res.value.key;
          try {
            jsonRes.Record = JSON.parse(res.value.value.toString("utf8"));
          } catch (err) {
            console.log(err);
            jsonRes.Record = res.value.value.toString("utf8");
          }
        }
        allResults.push(jsonRes);
      }
      if (res.done) {
        console.log("end of data");
        await iterator.close();
        console.info(allResults);
        return allResults;
      }
    }
  }
}

module.exports = IOTContract;
