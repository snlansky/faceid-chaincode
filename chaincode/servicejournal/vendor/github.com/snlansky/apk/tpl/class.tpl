//This file is generate by scripts,don't edit it
import com.kingdee.kchain.fabric.sdk.chaincode.TransactionResponse;
import com.kingdee.kchain.fabric.sdk.chaincode.payload.GsonPayloadParser;
import net.sf.json.JSONObject;
{{range $i, $_ := .Imports}}
{{$i}};
{{end}}

public class {{.ClassName}} extends Adapter{
    {{range $_, $field := .Fields}}
    {{$field}};
    {{end}}

    // public {{.ClassName}}(){}

    {{range $_, $method := .Methods}}
    {{$method.GetSig}}{
        String method="{{$.ClassName | $method.GenInterface}}";
        Object[] ojbs = {{$method.GetParams}};
        JSONObject result = JSONObject.fromObject(ojbs);
        GsonPayloadParser<{{$method.RetType}}> payload = GsonPayloadParser.newParser({{$method.RetType}}.class);
        TransactionResponse<{{$method.RetType}}> invoke = getChaincodeTemplate().prepareInvoke(method, payload)
                .user(getUser())
                .transientMap(new HashMap<String, byte[]>())
                .args(result.toString())
                .invoke();
        return invoke.getPayload();
    }
    {{end}}
}
