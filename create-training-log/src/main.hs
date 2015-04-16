import System.IO
import System.Locale as Locale
import Data.Time
import Data.Time.Format
import Data.Data
import qualified Data.Yaml as Yaml
import qualified Data.Aeson.TH as Aeson

data Workout = Workout { name :: String,
                         weight :: String,
                         sets :: String,
                         reps :: String,
                         exertion :: String } deriving (Eq, Show)

data TrainingLog = TrainingLog { date :: String,
                                 time :: String,
                                 duration :: String,
                                 bodyweight :: String,
                                 workout :: [Workout] } deriving (Eq, Show)

Aeson.deriveJSON Aeson.defaultOptions ''Workout
Aeson.deriveJSON Aeson.defaultOptions ''TrainingLog

main = do
  hSetBuffering stdout NoBuffering
  putStr "Enter date (YYYY-MM-DD) or [Enter] for today: "
  dateStr <- getLine
  putStr "Enter time: "
  timeStr <- getLine
  let
      dateTimeStr = dateStr ++ " " ++ timeStr
      dateTime = parseTime Locale.defaultTimeLocale "%F %l:%M%p" dateTimeStr :: Maybe UTCTime
  print dateTime
  putStr "Enter duration (e.g 1h30m): "
  durationStr <- getLine
  putStr "Enter bodyweight (e.g 70.0kg): "
  bodyweightStr <- getLine
  -- putStr "Enter Workout: "
  let
      trainingLog = TrainingLog{date=dateStr,
                                time=timeStr,
                                duration=durationStr,
                                bodyweight=bodyweightStr,
                                workout=[]}
      filename = "./" ++ dateStr ++ ".yml"
  print trainingLog
  Yaml.encodeFile filename trainingLog
