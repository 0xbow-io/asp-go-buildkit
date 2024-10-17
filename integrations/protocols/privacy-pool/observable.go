package privacypool

import (
	"math/big"

	watcher "github.com/0xBow-io/asp-go-buildkit/core/watcher"
	"github.com/0xBow-io/asp-go-buildkit/internal/erpc"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	logging "github.com/ipfs/go-log/v2"
	"github.com/cockroachdb/errors"
)

var (
	ErrorScopeNotFound    = errors.New("invalid scope")
	ErrorInstanceNotFound = errors.New("invalid instance")
	ErrorIterator         = errors.New("failed to create iterator")

	log = logging.Logger("watcher")
)

type observable uint

const (
	SEPOLIA_ETH_POOL_1 observable = iota
	SEPOLIA_ETH_POOL_2
	GNOSIS_XDAI_POOL_1
	GNSOSI_XDA_POOL_2
)

var _ watcher.Observable = (*observable)(nil)

// Observables returns the list of observables
func Observables() []watcher.Observable {
	return []watcher.Observable{
		SEPOLIA_ETH_POOL_1,
		SEPOLIA_ETH_POOL_2,
		GNOSIS_XDAI_POOL_1,
		GNSOSI_XDA_POOL_2,
	}
}

func (ob observable) ID() string {
	return [...]string{"SEPOLIA_ETH_POOL_1", "SEPOLIA_ETH_POOL_2", "GNOSIS_XDAI_POOL_1", "GNOSIS_XDAI_POOL_2"}[ob]
}

func strToBigInt(str string) *big.Int {
	i, _ := new(big.Int).SetString(str, 10)
	return i
}

func (ob observable) Scope() []byte {
	return [...][]byte{
		strToBigInt("15365509683721112532018974415132282847207162026665662018590046777583916671872").Bytes(),
		strToBigInt("1594601211935923806427821481643004967624986397998197460555337643549018639657").Bytes(),
		strToBigInt("11049869816642268564454296009173568684966369147224378104485796423384633924130").Bytes(),
		strToBigInt("19420586229045152356890556789607410844693215030122143238126523862419003191309").Bytes(),
	}[ob]
}

func (ob observable) ChainID() int {
	return [...]int{11155111, 11155111, 100, 100}[ob]
}

func (ob observable) Genesis() uint64 {
	return [...]uint64{6313019, 6454920, 34972988, 35827812}[ob]
}

func (ob observable) Address() common.Address {
	return [...]common.Address{
		common.HexToAddress("0x35F9acbaD838b12AA130Ef6386C14d847bdC1642"),
		common.HexToAddress("0x0C606138Aa02600c55e0d427cf4B2a7319a808fe"),
		common.HexToAddress("0x0C606138Aa02600c55e0d427cf4B2a7319a808fe"),
		common.HexToAddress("0x555eb8F3C1C2bEDa8e8eA69F8c51317470Ab8fC1"),
	}[ob]
}

func (ob observable) instance(adapter erpc.Backend) (*PrivacyPool, error) {
	return NewPrivacyPool(ob.Address(), adapter)
}

// Play returns a channel which streams the state of the observable
// starting from a given block number to a desired block number.
// State is derived from a state-transition event (`Record` events)
// and is published in the form of a serialized byte array.
func (ob observable) Play(adapter erpc.Backend, opts *bind.FilterOpts) (<-chan []byte, error) {
	var (
		instance, err = ob.instance(adapter)
		sink          = make(chan []byte, 24)
	)

	if err != nil || instance == nil {
		return nil, errors.Wrap(err, ErrorInstanceNotFound.Error())
	}

	iterator, err := instance.FilterRecord(opts)
	if err != nil {
		return nil, errors.Wrap(err, ErrorIterator.Error())
	}

	go func() {
		defer close(sink)
		defer iterator.Close()

		for {
			if record := iterator.Event; record != nil {
				sink <- new(State).DeriveFrom(ob.Scope(), record).Serialize()
			}
			ok := iterator.Next()
			if !ok {
				if err := iterator.Error(); err != nil {
					log.Errorw("privacypool/observable/Play: caught iterator error", "id", ob.ID(), "error", err)
					sink <- nil
					return
				}
				log.Debugw("privacypool/observable/Play: iterator done", "id", ob.ID())
				return
			}
		}
	}()

	return sink, nil
}

func (ob observable) Deserialize(bin []byte) watcher.State { return new(State).Deserialize(bin) }
